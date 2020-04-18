package keywords

import "github.com/jinzhu/gorm"

func NewKeywordService(db *gorm.DB) KeywordService {
	kg := &keyGorm{db: db}
	kv := newKeywordValidator(kg)
	return &keywordService{kv}
}

type KeywordService interface {
	KeywordDB
}

type keywordService struct {
	KeywordDB
}
