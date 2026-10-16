package extractor

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archives"
)

// extractArchive extracts a supported archive (zip, tar.xz, etc.) into destination.
// All writes go through an os.Root so entries can't escape the destination directory.
func (e *Extractor) extractArchive(ctx context.Context, filePath string, destination string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	format, stream, err := archives.Identify(ctx, filePath, file)
	if err != nil {
		return fmt.Errorf("identifying archive format: %w", err)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unsupported archive format: %s", format.Extension())
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	root, err := os.OpenRoot(destination)
	if err != nil {
		return err
	}
	defer root.Close()

	return extractor.Extract(ctx, stream, func(ctx context.Context, info archives.FileInfo) error {
		return e.writeEntry(root, info)
	})
}

func (e *Extractor) writeEntry(root *os.Root, info archives.FileInfo) error {
	name, err := localPath(info.NameInArchive)
	if err != nil {
		return err
	}

	if name == "." {
		return nil
	}

	if err := root.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	if hdr, ok := info.Header.(*tar.Header); ok && hdr.Typeflag == tar.TypeLink {
		target, err := localPath(info.LinkTarget)
		if err != nil {
			return err
		}

		return root.Link(target, name)
	}

	switch mode := info.Mode(); {
	case mode.IsDir():
		return root.MkdirAll(name, 0755)
	case mode&fs.ModeSymlink != 0:
		// Symlink targets aren't checked by os.Root, so make sure they stay inside the destination.
		target := filepath.FromSlash(info.LinkTarget)
		if filepath.IsAbs(target) || !filepath.IsLocal(filepath.Join(filepath.Dir(name), target)) {
			return fmt.Errorf("illegal symlink target %q for %q", info.LinkTarget, info.NameInArchive)
		}

		return root.Symlink(target, name)
	case mode.IsRegular():
		return writeFile(root, name, info)
	default:
		e.logger.Debug("skipping unsupported archive entry", map[string]interface{}{
			"name": info.NameInArchive,
			"mode": mode.String(),
		})
		return nil
	}
}

func writeFile(root *os.Root, name string, info archives.FileInfo) error {
	src, err := info.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	perm := info.Mode().Perm()
	if perm == 0 {
		perm = 0644
	}

	dst, err := root.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return err
	}

	return dst.Close()
}

// localPath converts an archive entry name into a relative path, rejecting anything that would escape the root.
func localPath(name string) (string, error) {
	cleaned := filepath.FromSlash(path.Clean(filepath.ToSlash(name)))
	if !filepath.IsLocal(cleaned) {
		return "", fmt.Errorf("illegal file path in archive: %q", name)
	}

	return cleaned, nil
}
