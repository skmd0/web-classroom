package keywords

import (
	"github.com/jinzhu/gorm"
	"wiki/models"
)

type Keyword struct {
	gorm.Model
	PostID      uint   `gorm:"not_null;index"`
	Title       string `gorm:"not_null"`
	Description string `gorm:"not_null"`
}

type KeywordDB interface {
	ByID(id uint) (*Keyword, error)
	ByPostID(id uint) (*[]Keyword, error)
	Create(keyword *Keyword) error
	Update(keyword *Keyword) error
	Delete(id uint) error
}

type keyGorm struct {
	db *gorm.DB
}

var _ KeywordDB = &keyGorm{}

func (kg *keyGorm) Create(key *Keyword) error {
	return kg.db.Create(key).Error
}

func (kg *keyGorm) Update(key *Keyword) error {
	return kg.db.Save(key).Error
}

func (kg *keyGorm) Delete(id uint) error {
	key := Keyword{Model: gorm.Model{ID: id}}
	return kg.db.Delete(&key).Error
}

// ByID looks up the user by the provided ID.
func (kg *keyGorm) ByID(id uint) (*Keyword, error) {
	var key Keyword
	db := kg.db.Where("id = ?", id)
	err := first(db, &key)
	if err != nil {
		return nil, err
	}
	return &key, err
}

func (kg *keyGorm) ByPostID(id uint) (*[]Keyword, error) {
	var keywords []Keyword
	err := kg.db.Where("post_id = ?", id).Find(&keywords).Error
	if err != nil {
		return nil, err
	}
	return &keywords, nil
}

// first executes a query from gorm.DB and writes data to dst by reference.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return models.ErrNotFound
	}
	return err
}
