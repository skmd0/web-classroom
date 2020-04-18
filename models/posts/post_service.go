package posts

import (
	"github.com/jinzhu/gorm"
	"wiki/models/keywords"
)

func NewPostService(db *gorm.DB, ks keywords.KeywordService) PostService {
	pg := &postGorm{db: db}
	uv := newPostValidator(pg)
	return &postService{
		PostDB: uv,
		ks:     ks,
	}
}

// PostService is a contract for the consumers
type PostService interface {
	PostDB
	CreatePost(*Post, *[]keywords.Keyword) error
	UpdatePost(*Post, *[]keywords.Keyword) error
	GetKeywords(uint) (*[]keywords.Keyword, error)
}

// Implementation of the PostService
type postService struct {
	PostDB
	ks keywords.KeywordService
}

func (ps *postService) CreatePost(post *Post, keys *[]keywords.Keyword) error {
	err := ps.PostDB.Create(post)
	if err != nil {
		return err
	}

	for _, key := range *keys {
		// ID of post is set after it is written in DB
		key.PostID = post.ID
		err = ps.ks.Create(&key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *postService) UpdatePost(post *Post, keys *[]keywords.Keyword) error {
	err := ps.PostDB.Update(post)
	if err != nil {
		return err
	}

	for _, key := range *keys {
		// ID of post is set after it is written in DB
		key.PostID = post.ID
		err = ps.ks.Update(&key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *postService) GetKeywords(postID uint) (*[]keywords.Keyword, error) {
	return ps.ks.ByPostID(postID)
}
