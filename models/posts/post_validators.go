package posts

import "wiki/models"

func newPostValidator(pdb PostDB) *postValidator {
	return &postValidator{
		PostDB: pdb,
	}
}

type postValidator struct {
	PostDB
}

// validator function type signature
type postValFunc func(*Post) error

// runPostValFunc loops through all the validator and returns err if any fail
func runPostValFunc(post *Post, fns ...postValFunc) error {
	for _, fn := range fns {
		if err := fn(post); err != nil {
			return err
		}
	}
	return nil
}

func (pv *postValidator) Create(post *Post) error {
	err := runPostValFunc(post,
		pv.titleRequired,
		pv.contentRequired,
		pv.userIDRequired)
	if err != nil {
		return err
	}
	return pv.PostDB.Create(post)
}

func (pv *postValidator) Update(post *Post) error {
	err := runPostValFunc(post,
		pv.titleRequired,
		pv.contentRequired,
		pv.userIDRequired)
	if err != nil {
		return err
	}
	return pv.PostDB.Update(post)
}

func (pv *postValidator) Delete(id uint) error {
	var post Post
	post.ID = id
	err := runPostValFunc(&post,
		pv.idGreaterThanZero)
	if err != nil {
		return err
	}
	return pv.PostDB.Delete(id)
}

func (pv *postValidator) idGreaterThanZero(post *Post) error {
	if post.ID == 0 {
		return models.ErrIdInvalid
	}
	return nil
}

func (pv *postValidator) userIDRequired(post *Post) error {
	if post.UserID <= 0 {
		return models.ErrUserIdRequired
	}
	return nil
}

func (pv *postValidator) titleRequired(post *Post) error {
	if post.Title == "" {
		return models.ErrTitleRequired
	}
	return nil
}

func (pv *postValidator) contentRequired(post *Post) error {
	if post.Content == "" {
		return models.ErrContentRequired
	}
	return nil
}
