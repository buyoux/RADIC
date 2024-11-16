package test

import (
	"RADIC/util"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

const P = 1000

var conMp = util.NewConcurrentHashMap(8, P)
var synMp = sync.Map{}

func readConMap() {
	for i := 0; i < P; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMp.Get(key)
	}
}

func writeConMap() {
	for i := 0; i < P; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		conMp.Set(key, 1)
	}
}

func readSynMap() {
	for i := 0; i < P; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMp.Load(key)
	}
}

func writeSynMap() {
	for i := 0; i < P; i++ {
		key := strconv.Itoa(int(rand.Int63()))
		synMp.Store(key, 1)
	}
}

func BenchmarkConMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const Q = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * Q)
		for i := 0; i < Q; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					readConMap()
				}
			}()
		}
		for i := 0; i < Q; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					writeConMap()
				}
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSynMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const Q = 300
		wg := sync.WaitGroup{}
		wg.Add(2 * Q)
		for i := 0; i < Q; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					readSynMap()
				}
			}()
		}
		for i := 0; i < Q; i++ {
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					writeSynMap()
				}
			}()
		}
		wg.Wait()
	}
}

// go test .\util\test\ -bench=Map -run=^$ -count=1 -benchmem -benchtime=3s
/*
goos: windows
goarch: amd64
pkg: gift/util/test
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkConMap-8              2        1737533100 ns/op        752222952 B/op  18099820 allocs/op
BenchmarkSynMap-8              1        4932080600 ns/op        432311992 B/op  12122118 allocs/op
PASS
ok      gift/util/test  10.687s
*/
