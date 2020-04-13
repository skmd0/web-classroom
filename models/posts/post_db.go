package posts

import "github.com/jinzhu/gorm"

// Post is representation of the post DB table
type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not_null;index"`
	Title   string `gorm:"not_null"`
	Content string `gorm:"not_null"`
}
