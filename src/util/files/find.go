package files

import (
	"io/fs"
	"path/filepath"
)

func FindSounds(ext string) []string {
	var a []string
	filepath.WalkDir("../sounds", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, d.Name())
		}
		return nil
	})
	return a
}
