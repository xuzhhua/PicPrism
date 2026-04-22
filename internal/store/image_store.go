package store

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ImageStore struct {
	db *sqlx.DB
}

func NewImageStore(db *sqlx.DB) *ImageStore {
	return &ImageStore{db: db}
}

func (s *ImageStore) Create(img *Image) error {
	_, err := s.db.NamedExec(`
		INSERT INTO images (id, filename, ext, size, width, height, mime_type, hash, created_at)
		VALUES (:id, :filename, :ext, :size, :width, :height, :mime_type, :hash, :created_at)
	`, img)
	return err
}

func (s *ImageStore) GetByID(id string) (*Image, error) {
	var img Image
	err := s.db.Get(&img, `SELECT * FROM images WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (s *ImageStore) GetByHash(hash string) (*Image, error) {
	var img Image
	err := s.db.Get(&img, `SELECT * FROM images WHERE hash = ?`, hash)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (s *ImageStore) List(f ListFilter) ([]Image, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 200 {
		f.Limit = 40
	}

	orderBy := "i.created_at DESC"
	switch f.Sort {
	case "oldest":
		orderBy = "i.created_at ASC"
	case "name":
		orderBy = "i.filename ASC"
	case "size":
		orderBy = "i.size DESC"
	}

	args := []interface{}{}
	where := []string{}

	if f.Tag != "" {
		where = append(where, `i.id IN (
			SELECT it.image_id FROM image_tags it
			JOIN tags t ON t.id = it.tag_id WHERE t.name = ?
		)`)
		args = append(args, f.Tag)
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// 总数
	var total int
	countArgs := append([]interface{}{}, args...)
	if err := s.db.QueryRow(
		fmt.Sprintf(`SELECT COUNT(*) FROM images i %s`, whereClause),
		countArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	listArgs := append(args, f.Limit, offset)
	rows, err := s.db.Queryx(
		fmt.Sprintf(`SELECT i.* FROM images i %s ORDER BY %s LIMIT ? OFFSET ?`, whereClause, orderBy),
		listArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.StructScan(&img); err != nil {
			return nil, 0, err
		}
		images = append(images, img)
	}
	return images, total, nil
}

func (s *ImageStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM images WHERE id = ?`, id)
	return err
}
