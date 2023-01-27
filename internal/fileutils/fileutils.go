package fileutils

import "os"

// WithFileOpened runs cb with fpath file opened, it doesn't check if fp.Close was successful
func WithFileOpened(fpath string, cb func(fp *os.File) error, flag int, mode os.FileMode) error {
	fp, err := os.OpenFile(fpath, flag, mode)
	if err != nil {
		return err
	}
	defer func() { _ = fp.Close() }()

	return cb(fp)
}
