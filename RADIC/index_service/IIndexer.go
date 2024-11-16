package indexservice

import "RADIC/types"

type IIndexer interface {
	AddDoc(doc types.Document) (int, error)
	DeleteDoc(docId string) int
	Search(query *types.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) []*types.Document
	Close() error
}
