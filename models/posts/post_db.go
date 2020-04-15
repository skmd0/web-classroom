package posts

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"wiki/models"
)

// Post is representation of the post DB table
type Post struct {
	gorm.Model
	UserID      uint   `gorm:"not_null;index"`
	Title       string `gorm:"not_null"`
	Content     string `gorm:"not_null"`
	ContentHTML template.HTML
}

type PostDB interface {
	ByID(id uint) (*Post, error)
	ByUserID(id uint) (*[]Post, error)
	ByUserIdWithLimit(id uint, limit int) (*[]Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uint) error
}

type postGorm struct {
	db *gorm.DB
}

var _ PostDB = &postGorm{}

func (pg *postGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}

func (pg *postGorm) Update(post *Post) error {
	return pg.db.Save(post).Error
}

func (pg *postGorm) Delete(id uint) error {
	post := Post{Model: gorm.Model{ID: id}}
	return pg.db.Delete(&post).Error
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

func (pg *postGorm) ByUserID(id uint) (*[]Post, error) {
	var posts []Post
	pg.db.Where("user_id = ?", id).Find(&posts)
	return &posts, nil
}

func (pg *postGorm) ByUserIdWithLimit(id uint, limit int) (*[]Post, error) {
	var posts []Post
	pg.db.Where("user_id = ?", id).Limit(limit).Find(&posts)
	return &posts, nil
}

// first executes a query from gorm.DB and writes data to dst by reference.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	return err
}
