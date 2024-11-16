package test

import (
	reverseindex "RADIC/internal/reverse_index"
	"fmt"
	"testing"

	"github.com/huandu/skiplist"
)

func TestIntersectionOfSkipList(t *testing.T) {
	l1 := skiplist.New(skiplist.Uint64)
	l1.Set(uint64(5), 0)
	l1.Set(uint64(1), 0)
	l1.Set(uint64(4), 0)
	l1.Set(uint64(9), 0)
	l1.Set(uint64(11), 0)
	l1.Set(uint64(7), 0)
	// skiplist内部自动排序，排序后为1，4，5，7，9，11

	l2 := skiplist.New(skiplist.Uint64)
	l2.Set(uint64(4), 0)
	l2.Set(uint64(5), 0)
	l2.Set(uint64(9), 0)
	l2.Set(uint64(8), 0)
	l2.Set(uint64(2), 0)

	l3 := skiplist.New(skiplist.Uint64)
	l3.Set(uint64(3), 0)
	l3.Set(uint64(5), 0)
	l3.Set(uint64(7), 0)
	l3.Set(uint64(9), 0)

	fmt.Println(" ")
	interset := reverseindex.IntersectionOfSkipList()
	if interset != nil {
		node := interset.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1")
	interset = reverseindex.IntersectionOfSkipList(l1)
	if interset != nil {
		node := interset.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1, l2")
	interset = reverseindex.IntersectionOfSkipList(l1, l2)
	if interset != nil {
		node := interset.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1, l2, l3")
	interset = reverseindex.IntersectionOfSkipList(l1, l2, l3)
	if interset != nil {
		node := interset.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}
}

func TestUnionsectionOfSkipList(t *testing.T) {
	l1 := skiplist.New(skiplist.Uint64)
	l1.Set(uint64(5), 0)
	l1.Set(uint64(1), 0)
	l1.Set(uint64(4), 0)
	l1.Set(uint64(9), 0)
	l1.Set(uint64(11), 0)
	l1.Set(uint64(7), 0)
	// skiplist内部自动排序，排序后为1，4，5，7，9，11

	l2 := skiplist.New(skiplist.Uint64)
	l2.Set(uint64(4), 0)
	l2.Set(uint64(5), 0)
	l2.Set(uint64(9), 0)
	l2.Set(uint64(8), 0)
	l2.Set(uint64(2), 0)
	l2.Set(uint64(10), 0)

	l3 := skiplist.New(skiplist.Uint64)
	l3.Set(uint64(3), 0)
	l3.Set(uint64(5), 0)
	l3.Set(uint64(6), 0)
	l3.Set(uint64(7), 0)
	l3.Set(uint64(9), 0)

	fmt.Println(" ")
	union := reverseindex.UnionsetOfSkipList()
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1")
	union = reverseindex.UnionsetOfSkipList(l1)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1, l2")
	union = reverseindex.UnionsetOfSkipList(l1, l2)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}

	fmt.Println("\nl1, l2, l3")
	union = reverseindex.UnionsetOfSkipList(l1, l2, l3)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key())
			node = node.Next()
		}
	}
}

// go test -v .\internal\reverse_index\test\  -run=^TestIntersectionOfSkipList$ -count=1
// go test -v .\internal\reverse_index\test\  -run=^TestUnionsectionOfSkipList$ -count=1
