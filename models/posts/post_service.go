package posts

import (
	"github.com/jinzhu/gorm"
)

func NewPostService(db *gorm.DB) PostService {
	pg := &postGorm{db: db}
	uv := newPostValidator(pg)
	return &postService{uv}
}

// PostService is a contract for the consumers
type PostService interface {
	PostDB
}

// Implementation of the PostService
type postService struct {
	PostDB
}
