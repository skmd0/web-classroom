package posts

import (
	"github.com/jinzhu/gorm"
	"wiki/models"
)

// Post is representation of the post DB table
type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not_null;index"`
	Title   string `gorm:"not_null"`
	Content string `gorm:"not_null"`
}

type PostDB interface {
	ByID(id uint) (*Post, error)
	Create(post *Post) error
}

type postGorm struct {
	db *gorm.DB
}

var _ PostDB = &postGorm{}

func (pg *postGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

// ByID looks up the user by the provided ID.
func (pg *postGorm) ByID(id uint) (*Post, error) {
	var post Post
	db := pg.db.Where("id = ?", id)
	err := first(db, &post)
	if err != nil {
		return nil, err
	}
	return &post, err
}

// first executes a query from gorm.DB and writes data to dst by reference.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	return err
}
