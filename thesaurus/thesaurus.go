package thesaurus

// Thesaurus 類語検索の汎用的なインターフェース
type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
