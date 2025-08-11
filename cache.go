package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("word not found")

type Cache struct {
	dir string
}

func NewCache() (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".cache", "anpi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Cache{dir: dir}, nil
}

func (c *Cache) path(word string) string {
	return filepath.Join(c.dir, word+".json")
}

func (c *Cache) notfoundPath(word string) string {
	return filepath.Join(c.dir, word+".notfound")
}

func (c *Cache) Read(word string) (*Entry, error) {
	if _, err := os.Stat(c.notfoundPath(word)); err == nil {
		return nil, ErrNotFound
	}
	p := c.path(word)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck
	var e Entry
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (c *Cache) Write(word string, e *Entry, found bool) error {
	if !found {
		f, err := os.Create(c.notfoundPath(word))
		if err != nil {
			return err
		}
		return f.Close()
	}
	p := c.path(word)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	return json.NewEncoder(f).Encode(e)
}

func (c *Cache) Clear(word string) error {
	_ = os.Remove(c.notfoundPath(word))
	return os.Remove(c.path(word))
}
