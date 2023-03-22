package main

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"time"
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

func Released(path string) ([]DirEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var released []DirEntry
	r := csv.NewReader(f)
	r.Comma = ' '
	r.FieldsPerRecord = 2
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		now := time.Now()
		t, err := time.ParseInLocation(
			"2006-01-02 15:04:05",
			record[0],
			now.Location(),
		)
		if err != nil {
			return nil, err
		}
		if t.After(now) {
			continue
		}

		name := record[1]
		if len(name) > 32 {
			name = name[:29] + "..."
		}
		released = append(released, DirEntry{
			Name: name,
			Path: record[1],
			Time: t.Format("Jan _2 15:04 2006"),
		})
	}

	return released, nil
}
