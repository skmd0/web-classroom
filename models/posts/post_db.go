package posts

import "github.com/jinzhu/gorm"

// Post is representation of the post DB table
type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not_null;index"`
	Title   string `gorm:"not_null"`
	Content string `gorm:"not_null"`
}

type PostDB interface {
	Create(post *Post) error
}

type postGorm struct {
	db *gorm.DB
}

var _ PostDB = &postGorm{}

func (pg *postGorm) Create(post *Post) error {
	return pg.db.Create(post).Error
}
