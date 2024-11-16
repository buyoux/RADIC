package reverseindex

import "RADIC/types"

type IReverseIndexer interface {
	Add(doc types.Document)
	Delete(IntId uint64, keyword *types.Keyword)
	Search(q *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []string
}
