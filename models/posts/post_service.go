package posts

import (
	"github.com/jinzhu/gorm"
)

func NewPostService(db *gorm.DB) PostService {
	pg := &postGorm{db: db}
	uv := newPostValidator(pg)
	return &postService{uv}
}

type PostService interface {
	PostDB
}

type postService struct {
	PostDB
}
