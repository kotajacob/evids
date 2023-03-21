package main

import (
	"os"
	"path/filepath"
)

type DirEntry struct {
	Name string
	Path string
	Size int64
	Time string
}

func ListDir(path string) ([]DirEntry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []DirEntry
	for _, e := range dirEntries {
		name := e.Name()
		if len(name) > 28 {
			name = name[:25] + "..."
		}
		if e.IsDir() {
			name += "/"
		}

		info, err := e.Info()
		if err != nil {
			return nil, err
		}

		entries = append(entries, DirEntry{
			Name: name,
			Path: filepath.Join(filepath.Base(path), e.Name()),
			Size: info.Size(),
			Time: info.ModTime().Format("Jan _2 15:04 2006"),
		})
	}
	return entries, nil
}
