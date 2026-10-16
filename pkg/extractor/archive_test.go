package extractor

import (
	"archive/tar"
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rocketblend/rocketblend/pkg/types"
)

type tarEntry struct {
	name     string
	body     string
	typeflag byte
	linkname string
}

func writeTar(t *testing.T, path string, entries []tarEntry) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tw := tar.NewWriter(f)
	for _, entry := range entries {
		hdr := &tar.Header{
			Name:     entry.name,
			Typeflag: entry.typeflag,
			Linkname: entry.linkname,
			Mode:     0755,
			Size:     int64(len(entry.body)),
		}
		if entry.typeflag != tar.TypeReg {
			hdr.Size = 0
		}

		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		if entry.typeflag == tar.TypeReg {
			if _, err := tw.Write([]byte(entry.body)); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeZip(t *testing.T, path string, files map[string]string) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}

	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
}

func assertFile(t *testing.T, path string, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Errorf("%s: got %q, want %q", path, got, want)
	}
}

func TestExtractZip(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "blender.zip")
	output := filepath.Join(dir, "out")

	writeZip(t, archive, map[string]string{
		"blender/blender.exe":      "binary",
		"blender/4.2/readme.txt":   "hello",
		"blender/4.2/scripts/a.py": "print()",
	})

	e, _ := New()
	if err := e.Extract(context.Background(), &types.ExtractOpts{Path: archive, OutputPath: output}); err != nil {
		t.Fatal(err)
	}

	assertFile(t, filepath.Join(output, "blender", "blender.exe"), "binary")
	assertFile(t, filepath.Join(output, "blender", "4.2", "readme.txt"), "hello")
	assertFile(t, filepath.Join(output, "blender", "4.2", "scripts", "a.py"), "print()")
}

func TestExtractTarLinks(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "blender.tar")
	output := filepath.Join(dir, "out")

	writeTar(t, archive, []tarEntry{
		{name: "blender/", typeflag: tar.TypeDir},
		{name: "blender/lib/libfoo.so.1", body: "lib", typeflag: tar.TypeReg},
		{name: "blender/lib/libfoo.so", typeflag: tar.TypeSymlink, linkname: "libfoo.so.1"},
		{name: "blender/lib/libfoo-hard.so", typeflag: tar.TypeLink, linkname: "blender/lib/libfoo.so.1"},
	})

	e, _ := New()
	if err := e.Extract(context.Background(), &types.ExtractOpts{Path: archive, OutputPath: output}); err != nil {
		t.Fatal(err)
	}

	assertFile(t, filepath.Join(output, "blender", "lib", "libfoo.so"), "lib")
	assertFile(t, filepath.Join(output, "blender", "lib", "libfoo-hard.so"), "lib")

	target, err := os.Readlink(filepath.Join(output, "blender", "lib", "libfoo.so"))
	if err != nil {
		t.Fatal(err)
	}

	if target != "libfoo.so.1" {
		t.Errorf("symlink target: got %q, want %q", target, "libfoo.so.1")
	}
}

func TestExtractRejectsTraversal(t *testing.T) {
	tests := []struct {
		name  string
		entry tarEntry
	}{
		{
			name:  "relative path escape",
			entry: tarEntry{name: "../evil.txt", body: "x", typeflag: tar.TypeReg},
		},
		{
			name:  "absolute path",
			entry: tarEntry{name: "/tmp/evil.txt", body: "x", typeflag: tar.TypeReg},
		},
		{
			name:  "symlink escape",
			entry: tarEntry{name: "link", typeflag: tar.TypeSymlink, linkname: "../../etc/passwd"},
		},
		{
			name:  "absolute symlink",
			entry: tarEntry{name: "link", typeflag: tar.TypeSymlink, linkname: "/etc/passwd"},
		},
		{
			name:  "hard link escape",
			entry: tarEntry{name: "link", typeflag: tar.TypeLink, linkname: "../outside.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			archive := filepath.Join(dir, "evil.tar")
			output := filepath.Join(dir, "out")

			writeTar(t, archive, []tarEntry{tt.entry})

			e, _ := New()
			if err := e.Extract(context.Background(), &types.ExtractOpts{Path: archive, OutputPath: output}); err == nil {
				t.Fatal("expected error, got nil")
			}

			if _, err := os.Stat(filepath.Join(dir, "evil.txt")); !os.IsNotExist(err) {
				t.Error("file was written outside the output directory")
			}
		})
	}
}
