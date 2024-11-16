package util

import (
	"sync"

	"golang.org/x/exp/maps"

	farmhash "github.com/leemcloughlin/gofarmhash"
)

type ConcurrentHashMap struct {
	mps   []map[string]any // 由多个小map构成
	seg   int              // 小map的个数
	locks []sync.RWMutex   // 每个小map配一把读写锁。避免全局只有一把锁，影响性能
	seed  uint32           // 每次执行farmhash传统一的seed
}

// cap预估大map中容纳多少元素，seg内部包含几个小map
func NewConcurrentHashMap(seg, cap int) *ConcurrentHashMap {
	mps := make([]map[string]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := 0; i < seg; i++ {
		mps[i] = make(map[string]any, cap/seg)
		locks[i] = sync.RWMutex{}
	}
	return &ConcurrentHashMap{
		mps:   mps,
		seg:   seg,
		locks: locks,
		seed:  0,
	}
}

// 判断key对应到哪个小map
func (m *ConcurrentHashMap) getSegIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), m.seed))
	return hash % m.seg
}

// 写入<key, value>
func (m *ConcurrentHashMap) Set(key string, value any) {
	// index代表大map中小map的下标
	index := m.getSegIndex(key)
	m.locks[index].Lock()
	defer m.locks[index].Unlock()
	m.mps[index][key] = value
}

// 根据key读取value
func (m *ConcurrentHashMap) Get(key string) (any, bool) {
	index := m.getSegIndex(key)
	m.locks[index].RLock()
	defer m.locks[index].RUnlock()
	value, exists := m.mps[index][key]
	return value, exists
}

func (m *ConcurrentHashMap) CreateIterator() *ConcurrentHashMapIterator {
	keys := make([][]string, 0, len(m.mps))
	for _, mp := range m.mps {
		row := maps.Keys(mp)
		keys = append(keys, row)
	}
	return &ConcurrentHashMapIterator{
		cm:       m,
		keys:     keys,
		rowIndex: 0,
		colIndex: 0,
	}
}

type MapEntry struct {
	key   string
	value any
}

// 迭代器Iterator模式
type MapIterator interface {
	Next() *MapEntry
}

type ConcurrentHashMapIterator struct {
	cm       *ConcurrentHashMap
	keys     [][]string
	rowIndex int
	colIndex int
}

func (iter *ConcurrentHashMapIterator) Next() *MapEntry {
	// rowIndex超过了key的数量
	if iter.rowIndex >= len(iter.keys) {
		return nil
	}
	row := iter.keys[iter.rowIndex]
	if len(row) == 0 { // 本行为空
		iter.rowIndex += 1
		return iter.Next() // 进入递归，因为下一行可能依然为空
	}
	// 根据下标访问切片元素时，一定注意不要出现切片越界，
	// 即使下标为0，切片为空时依然会出现切片下标越界异常， out of range
	key := row[iter.colIndex]
	value, _ := iter.cm.Get(key)
	if iter.colIndex >= len(row)-1 {
		iter.rowIndex += 1
		iter.colIndex = 0
	} else {
		iter.colIndex += 1
	}
	return &MapEntry{key: key, value: value}
}
