package keywords

func newKeywordValidator(kdb KeywordDB) *keywordValidator {
	return &keywordValidator{KeywordDB: kdb}
}

type keywordValidator struct {
	KeywordDB
}
