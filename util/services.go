package util

import (
	"github.com/jinzhu/gorm"
	"wiki/models/posts"
	"wiki/models/users"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		Post: posts.NewPostService(db),
		User: users.NewUserService(db),
		db:   db,
	}, nil
}

type Services struct {
	Post posts.PostService
	User users.UserService
	db   *gorm.DB
}

// Close closes the database connection.
func (s *Services) Close() error {
	return s.db.Close()
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&users.User{}, &posts.Post{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate tries to automatically migrate the DB schema changes
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&users.User{}, &posts.Post{}).Error
}
