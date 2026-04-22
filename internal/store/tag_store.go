package store

import (
	"github.com/jmoiron/sqlx"
)

type TagStore struct {
	db *sqlx.DB
}

func NewTagStore(db *sqlx.DB) *TagStore {
	return &TagStore{db: db}
}

// UpsertTags 插入不存在的标签，返回全部对应 Tag 记录
func (s *TagStore) UpsertTags(names []string) ([]Tag, error) {
	tags := make([]Tag, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		_, _ = s.db.Exec(`INSERT OR IGNORE INTO tags (name) VALUES (?)`, name)
		var t Tag
		if err := s.db.Get(&t, `SELECT * FROM tags WHERE name = ?`, name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// SetImageTags 全量替换某图片的标签
func (s *TagStore) SetImageTags(imageID string, tagIDs []int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM image_tags WHERE image_id = ?`, imageID); err != nil {
		return err
	}
	for _, tid := range tagIDs {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO image_tags (image_id, tag_id) VALUES (?, ?)`, imageID, tid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetImageTags 返回某图片的所有标签
func (s *TagStore) GetImageTags(imageID string) ([]Tag, error) {
	var tags []Tag
	err := s.db.Select(&tags, `
		SELECT t.id, t.name FROM tags t
		JOIN image_tags it ON it.tag_id = t.id
		WHERE it.image_id = ?
		ORDER BY t.name
	`, imageID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// ListAll 返回全部标签及各自图片数量
func (s *TagStore) ListAll() ([]TagWithCount, error) {
	var tags []TagWithCount
	err := s.db.Select(&tags, `
		SELECT t.id, t.name, COUNT(it.image_id) AS count
		FROM tags t
		LEFT JOIN image_tags it ON it.tag_id = t.id
		GROUP BY t.id
		ORDER BY count DESC, t.name ASC
	`)
	return tags, err
}

// DeleteUnused 删除没有关联图片的标签
func (s *TagStore) DeleteUnused() error {
	_, err := s.db.Exec(`
		DELETE FROM tags WHERE id NOT IN (SELECT DISTINCT tag_id FROM image_tags)
	`)
	return err
}
