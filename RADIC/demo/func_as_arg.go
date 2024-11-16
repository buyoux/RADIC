package course

import (
	"container/ring"
	"fmt"
	"slices"
)

// 函数作为参数

var a func(int, string)

func TraverseRing(ring *ring.Ring) {
	// 通过Do()来遍历环Ring。函数参数用来指定如何处理Ring中的元素
	ring.Do(func(i interface{}) {
		fmt.Printf("%v", i)
	})
	fmt.Println()
}

func compare1(a, b *Doc) int {
	return a.Id - b.Id
}

func compare2(a, b *Doc) int {
	return b.Id - a.Id
}

func SortDoc1(docs []*Doc, compare func(a, b *Doc) int) {
	slices.SortFunc(docs, compare)
}

type IDocComparator interface {
	Compare(a, b *Doc) int
}

type PositiveOrder struct{}

func (PositiveOrder) Compare(a, b *Doc) int {
	return a.Id - b.Id
}

type ReversedOrder struct{}

func (ReversedOrder) Compare(a, b *Doc) int {
	return b.Id - a.Id
}

func SortDoc2(docs []*Doc, comparator IDocComparator) {
	slices.SortFunc(docs, comparator.Compare)
}

func funcArgVsInterface() {
	docs := make([]*Doc, 10)
	// 通过指定函数来确定正序还是逆序，即compare1和compare2
	// 函数作为参数
	SortDoc1(docs, compare1)
	SortDoc1(docs, compare2)

	// 接口作为参数，通过接口实现的函数来规定排序
	SortDoc2(docs, PositiveOrder{})
	SortDoc2(docs, ReversedOrder{})
}
