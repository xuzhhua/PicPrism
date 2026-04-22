package store

import (
	"time"
)

type Image struct {
	ID        string    `db:"id"`
	Filename  string    `db:"filename"`
	Ext       string    `db:"ext"`
	Size      int64     `db:"size"`
	Width     int       `db:"width"`
	Height    int       `db:"height"`
	MimeType  string    `db:"mime_type"`
	Hash      string    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
	Tags      []Tag     `db:"-"`
}

type Tag struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type TagWithCount struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Count int    `db:"count"`
}

type ListFilter struct {
	Tag   string
	Page  int
	Limit int
	Sort  string // newest | oldest | name | size
}
