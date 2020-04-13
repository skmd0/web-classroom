package posts

func newPostValidator(pdb PostDB) *postValidator {
	return &postValidator{
		PostDB: pdb,
	}
}

type postValidator struct {
	PostDB
}
