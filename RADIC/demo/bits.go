package course

import (
	"fmt"
	"math/bits"
)

func IsBit1(n uint64, i int) bool {
	if i > 64 {
		panic(i)
	}
	var c uint64 = 1 << (i - 1)
	if n&c == c {
		return true
	} else {
		return false
	}
}

func SetBit1(n uint64, i int) uint64 {
	if i > 64 {
		panic(i)
	}
	var c uint64 = 1 << (i - 1)
	return n | c
}

func CountBit1(n uint64) int {
	c := uint64(1)
	sum := 0
	for i := 0; i < 64; i++ {
		if n&c == c {
			sum += 1
		}
		c <<= 1
	}
	return bits.OnesCount64(n)
	// return sum
}

const (
	// 三种写法 完全等价
	// MALE = 1<<0
	// VIP = 1<<1
	// WEEK_ACTIVE = 1<<2

	// MALE = 1<<iota
	// VIP = 1<<iota
	// WEEK_ACTIVE = 1<<iota

	MALE = 1 << iota
	VIP
	WEEK_ACTIVE
)

type Candidate struct {
	Id     int
	Gender string
	Vip    bool
	Active int
	Bits   uint64
}

func (c *Candidate) SetMale() {
	c.Gender = "男"
	c.Bits |= MALE
}

func (c *Candidate) SetVip() {
	c.Vip = true
	c.Bits |= VIP
}

func (c *Candidate) SetActive(day int) {
	c.Active = day
	if day <= 7 {
		c.Bits |= WEEK_ACTIVE
	}
}

func (c Candidate) NormalFilter(male, vip, weekActive bool) bool {
	if male && c.Gender != "男" {
		return false
	}
	if vip && !c.Vip {
		return false
	}
	if weekActive && c.Active > 7 {
		return false
	}
	return true
}

// on代表多个条件同时满足的bits
func (c Candidate) BitsFilter(on uint64) bool {
	return c.Bits&on == on
}

type BitMap struct {
	Table uint64 // bits
}

func CreateBitMap(min int, arr []int) *BitMap {
	// bitMap := new(BitMap)
	bitMap := &BitMap{}
	for _, ele := range arr {
		index := ele - min
		bitMap.Table = SetBit1(bitMap.Table, index)
	}
	return bitMap
}

// 位图求交集
func IntersectionOfBitMap(bm1, bm2 *BitMap, min int) []int {
	rect := make([]int, 0, 100)
	s := bm1.Table & bm2.Table
	fmt.Printf("\n%064b", s)
	for i := 1; i <= 64; i++ {
		if IsBit1(s, i) {
			rect = append(rect, i+min)
		}
	}
	return rect
}

func IntersectionOfOrderedList(arr, brr []int) []int {
	m, n := len(arr), len(brr)
	if m == 0 || n == 0 {
		return []int{}
	}
	rect := make([]int, 0, 100)
	var i, j int
	for i < m && j < n {
		if arr[i] == brr[j] {
			rect = append(rect, arr[i])
			i++
			j++
		} else if arr[i] < brr[j] {
			i++
		} else {
			j++
		}
	}
	return rect
}
