package reverseindex

import (
	"RADIC/types"
	"RADIC/util"
	"runtime"
	"sync"

	"github.com/huandu/skiplist"
	farmhash "github.com/leemcloughlin/gofarmhash"
)

// 倒排索引整体上是个map，map的value是一个list
type SkipListReverseIndex struct {
	table *util.ConcurrentHashMap // 分段map，并发安全
	locks []sync.RWMutex          // 修改倒排索引时，相同的key需要去竞争同一把锁
}

// DocNumEstimate是预估的doc数量,新建一个倒排索引
func NewSkipListReverseIndex(DocNumEstimate int) *SkipListReverseIndex {
	indexer := new(SkipListReverseIndex)
	indexer.table = util.NewConcurrentHashMap(runtime.NumCPU(), DocNumEstimate)
	indexer.locks = make([]sync.RWMutex, 1000)
	return indexer
}

func (indexer SkipListReverseIndex) getLock(key string) *sync.RWMutex {
	n := int(farmhash.Hash32WithSeed([]byte(key), 0))
	return &indexer.locks[n%len(indexer.locks)]
}

type SkipListValue struct {
	Id          string
	BitsFeature uint64
}

// 添加一个doc
func (indexer *SkipListReverseIndex) Add(doc types.Document) {
	for _, keyword := range doc.Keywords {
		key := keyword.ToString()
		lock := indexer.getLock(key)
		lock.Lock()
		sklValue := SkipListValue{doc.Id, doc.BitsFeature}
		if value, exists := indexer.table.Get(key); exists {
			// table(也就是ConcurrentMap)的value即这里的list，也是一个map即跳表
			list := value.(*skiplist.SkipList)
			// IntId作为SkipList的key，而value里则包含了业务侧的IntId,BitsFeature
			list.Set(doc.IntId, sklValue)
		} else {
			list := skiplist.New(skiplist.Uint64)
			list.Set(doc.IntId, sklValue)
			indexer.table.Set(key, list)
		}
		lock.Unlock()
	}
}

func (indexer *SkipListReverseIndex) Delete(IntId uint64, keyword *types.Keyword) {
	key := keyword.ToString()
	lock := indexer.getLock(key)
	lock.Lock()
	if value, exists := indexer.table.Get(key); exists {
		list := value.(*skiplist.SkipList)
		list.Remove(IntId)
	}
	lock.Unlock()
}

// 求多个跳表SkipList的交集
func IntersectionOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}
	result := skiplist.New(skiplist.Uint64)
	currNodes := make([]*skiplist.Element, len(lists))
	for i, list := range lists {
		if list == nil || list.Len() == 0 {
			return nil
		}
		currNodes[i] = list.Front()
	}
	for {
		maxList := make(map[int]struct{}, len(lists))
		var maxValue uint64 = 0
		for i, node := range currNodes {
			if node.Key().(uint64) > maxValue {
				maxValue = node.Key().(uint64)
				maxList = map[int]struct{}{i: {}}
			} else if node.Key().(uint64) == maxValue {
				maxList[i] = struct{}{}
			}
		}
		// 所有node的值都一样大，则找到一个交集
		if len(maxList) == len(currNodes) {
			result.Set(currNodes[0].Key(), currNodes[0].Value)
			// 所有node均需往后移
			for i, node := range currNodes {
				currNodes[i] = node.Next()
				if currNodes[i] == nil {
					return result
				}
			}
		} else {
			for i, node := range currNodes {
				// maxValue的跳表不动，比maxValue小的往后移
				if _, exists := maxList[i]; !exists {
					// 不能使用node=node.Next(),for range是取得的是值拷贝
					currNodes[i] = node.Next()
					// 其中有一条SkipList走完，说明不会有新的交集产生，即退出返回
					if currNodes[i] == nil {
						return result
					}
				}
			}
		}
	}
}

// 求多个跳表SkipList的并集
func UnionsetOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}
	result := skiplist.New(skiplist.Uint64)
	keySet := make(map[any]struct{}, 1000)
	for _, list := range lists {
		if list == nil {
			continue
		}
		node := list.Front()
		for node != nil {
			if _, exists := keySet[node.Key()]; !exists {
				result.Set(node.Key(), node.Value)
				keySet[node.Key()] = struct{}{}
			}
			node = node.Next()
		}
	}
	return result
}

// 按照bits特征进行过滤
func (indexer SkipListReverseIndex) FilterByBits(bits uint64, onFlag uint64,
	offFlag uint64, orFlags []uint64) bool {
	// onFlag条件所有bit必须全部命中
	if bits&onFlag != onFlag {
		return false
	}
	// offFlag条件所有bit必须全部不命中
	if bits&offFlag != 0 {
		return false
	}
	// 对于所有orFlag,每个orFlag至少有一个bit命中
	for _, orFlag := range orFlags {
		// 对于单个orFlag 只要有一个bit命中即可
		if orFlag > 0 && bits&orFlag <= 0 {
			return false
		}
	}
	return true
}

// 搜索，返回SkipList
func (indexer SkipListReverseIndex) search(query *types.TermQuery, onFlag uint64,
	offFlag uint64, orFlags []uint64) *skiplist.SkipList {
	// 叶子节点，TermQuery是一个Keyword
	if query.Keyword != nil {
		keyword := query.Keyword.ToString()
		if value, exists := indexer.table.Get(keyword); exists {
			// result开辟一个空的跳表作为返回值
			result := skiplist.New(skiplist.Uint64)
			list := value.(*skiplist.SkipList)
			// util.Log.Printf("")
			node := list.Front()
			for node != nil {
				intId := node.Key().(uint64)
				skv, _ := node.Value.(SkipListValue)
				flag := skv.BitsFeature
				// 确保有效元素intId都大于0
				if intId > 0 && indexer.FilterByBits(flag, onFlag, offFlag, orFlags) {
					result.Set(intId, skv)
				}
				// node采用迭代器模式
				node = node.Next()
			}
			return result
		}
	} else if len(query.Must) > 0 {
		results := make([]*skiplist.SkipList, 0, len(query.Must))
		for _, must := range query.Must {
			results = append(results, indexer.search(must, onFlag, offFlag, orFlags))
		}
		return IntersectionOfSkipList(results...)
	} else if len(query.Should) > 0 {
		results := make([]*skiplist.SkipList, 0, len(query.Must))
		for _, should := range query.Should {
			results = append(results, indexer.search(should, onFlag, offFlag, orFlags))
		}
		return UnionsetOfSkipList(results...)
	}
	return nil
}

// 面向业务侧调用方法，搜索返回docId
func (indexer SkipListReverseIndex) Search(query *types.TermQuery, onFlag uint64,
	offFlag uint64, orFlags []uint64) []string {
	// 调用search搜索query条件，返回结果为满足条件的一个跳表
	skipListResult := indexer.search(query, onFlag, offFlag, orFlags)
	if skipListResult == nil {
		return nil
	}
	docIdResult := make([]string, 0, skipListResult.Len())
	node := skipListResult.Front()
	for node != nil {
		skv, _ := node.Value.(SkipListValue)
		docIdResult = append(docIdResult, skv.Id)
		node = node.Next()
	}
	return docIdResult
}
