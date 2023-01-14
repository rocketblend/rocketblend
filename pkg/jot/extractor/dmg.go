package extractor

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type DMGReader struct {
	r       io.ReadSeekCloser
	Entries []*DMGFileEntry
}

type DMGFileEntry struct {
	name   string
	isDir  bool
	mode   os.FileMode
	start  int64
	length int64
	r      io.ReadSeeker
}

func (e *DMGFileEntry) Name() string {
	return e.name
}

func (e *DMGFileEntry) IsDir() bool {
	return e.isDir
}

func (e *DMGFileEntry) Mode() os.FileMode {
	return e.mode
}

func (e *DMGFileEntry) Read(p []byte) (n int, err error) {
	return io.ReadFull(io.LimitReader(e.r, e.length), p)
}

func (r *DMGReader) Close() error {
	return r.r.Close()
}

func NewDMGReader(r io.ReadSeekCloser) (*DMGReader, error) {
	dmgReader := &DMGReader{r: r}
	err := dmgReader.readEntries()
	if err != nil {
		return nil, err
	}
	return dmgReader, nil
}

func (r *DMGReader) readEntries() error {
	var offset int64
	for {
		// Read the entry's header
		header := make([]byte, 12)
		_, err := io.ReadFull(r.r, header)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Check if the entry is a directory or a file
		isDir := header[11] == 1
		entryLength := binary.BigEndian.Uint32(header[:4])
		entryNameLength := binary.BigEndian.Uint32(header[4:8])
		entryName := make([]byte, entryNameLength)
		_, err = io.ReadFull(r.r, entryName)
		if err != nil {
			return err
		}
		entryName = bytes.Trim(entryName, "\x00")

		var mode os.FileMode
		if isDir {
			mode = os.ModeDir | 0755
		} else {
			mode = 0644
		}

		// Create a new DMGFileEntry and add it to the entries slice
		entry := &DMGFileEntry{
			name:   string(entryName),
			isDir:  isDir,
			mode:   mode,
			start:  offset + 12 + int64(entryNameLength),
			length: int64(entryLength) - int64(entryNameLength) - 12,
			r:      r.r,
		}
		r.Entries = append(r.Entries, entry)

		// Update the offset for the next entry
		offset += int64(entryLength)
	}
	return nil
}
