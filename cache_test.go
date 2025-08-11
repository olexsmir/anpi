package main

import (
	"errors"
	"testing"

	"olexsmir.xyz/x/is"
)

func testCache(t *testing.T) *Cache {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	c, err := NewCache()
	is.Err(t, err, nil)
	return c
}

func testEntry() *Entry {
	return &Entry{
		Word: "test",
		POSBlocks: []POSBlock{
			{POS: "noun", IPA: "/test/", Senses: []Sense{{Definition: "a test"}}},
		},
	}
}

func TestCache_WriteRead(t *testing.T) {
	c := testCache(t)
	entry := testEntry()
	is.Err(t, c.Write("test", entry, true), nil)

	got, err := c.Read("test")
	is.Err(t, err, nil)
	is.Equal(t, entry.Word, got.Word)
	is.Equal(t, len(entry.POSBlocks), len(got.POSBlocks))
	is.Equal(t, entry.POSBlocks[0].POS, got.POSBlocks[0].POS)
	is.Equal(t, entry.POSBlocks[0].IPA, got.POSBlocks[0].IPA)
	is.Equal(t, entry.POSBlocks[0].Senses[0].Definition, got.POSBlocks[0].Senses[0].Definition)
}

func TestCache_Read_miss(t *testing.T) {
	c := testCache(t)
	_, err := c.Read("nonexistent")
	is.Err(t, err, "no such file or directory")
}

func TestCache_Clear(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("test", testEntry(), true), nil)
	is.Err(t, c.Clear("test"), nil)
	_, err := c.Read("test")
	is.Err(t, err, "no such file or directory")
}

func TestCache_Clear_notFound(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("test", nil, false), nil)
	is.Err(t, c.Clear("test"), "no such file or directory")
	_, err := c.Read("test")
	is.Err(t, err, "no such file or directory")
}

func TestCache_Clear_missing(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Clear("nonexistent"), "no such file or directory")
}

func TestCache_Read_isolationBetweenWords(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("foo", testEntry(), true), nil)
	_, err := c.Read("bar")
	is.Err(t, err, "no such file or directory")
}

func TestCache_Write_notFound(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("xyz", nil, false), nil)
	_, err := c.Read("xyz")
	is.Equal(t, true, errors.Is(err, ErrNotFound))
}

func TestCache_Write_notFound_overridesEntry(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("test", testEntry(), true), nil)
	is.Err(t, c.Write("test", nil, false), nil)
	_, err := c.Read("test")
	is.Equal(t, true, errors.Is(err, ErrNotFound))
}

func TestCache_Write_notFound_isolation(t *testing.T) {
	c := testCache(t)
	is.Err(t, c.Write("foo", nil, false), nil)
	is.Err(t, c.Write("bar", testEntry(), true), nil)

	_, err := c.Read("foo")
	is.Equal(t, true, errors.Is(err, ErrNotFound))

	got, err := c.Read("bar")
	is.Err(t, err, nil)
	is.Equal(t, "test", got.Word)
}
