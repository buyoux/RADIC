package kvdb

import (
	"errors"
	"sync/atomic"

	bolt "go.etcd.io/bbolt"
)

// kv不存在也被认为是一种error
var NoDataErr = errors.New("No Data")

// Bolt blot store struct
type Bolt struct {
	db     *bolt.DB
	path   string
	bucket []byte
}

// Builder生成器模式
func (s *Bolt) WithDataPath(path string) *Bolt {
	s.path = path
	return s
}

func (s *Bolt) WithBucket(bucket string) *Bolt {
	s.bucket = []byte(bucket)
	return s
}

// OpenBolt open Blot store
func (s *Bolt) Open() error {
	DataDir := s.GetDbPath()
	db, err := bolt.Open(DataDir, 0o600, bolt.DefaultOptions)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.bucket)
		return err
	})
	if err != nil {
		db.Close()
		return err
	} else {
		s.db = db
		return nil
	}
}

func (s *Bolt) GetDbPath() string {
	return s.path
}

func (s *Bolt) Set(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Put(key, value)
	})
}

func (s *Bolt) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("key value not the same length!")
	}
	var err error
	err = s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			value := values[i]
			tx.Bucket(s.bucket).Put(key, value)
		}
		return nil
	})
	return err
}

func (s *Bolt) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(s.bucket).Get(key)
		return nil
	})
	if len(value) == 0 {
		return nil, NoDataErr
	}
	return value, err
}

func (s *Bolt) BatchGet(keys [][]byte) ([][]byte, error) {
	var err error
	values := make([][]byte, len(keys))
	err = s.db.Batch(func(tx *bolt.Tx) error {
		for i, key := range keys {
			value := tx.Bucket(s.bucket).Get(key)
			values[i] = value
		}
		return nil
	})
	if len(values) == 0 {
		return nil, NoDataErr
	}
	return values, err
}

func (s *Bolt) Delete(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(s.bucket).Delete(key)
	})
}

func (s *Bolt) BatchDelete(keys [][]byte) error {
	var err error
	err = s.db.Batch(func(tx *bolt.Tx) error {
		for _, key := range keys {
			tx.Bucket(s.bucket).Delete(key)
		}
		return nil
	})
	return err
}

func (s *Bolt) Has(key []byte) bool {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(s.bucket).Get(key)
		return nil
	})
	if len(value) == 0 || err != nil {
		return false
	}
	return true
}

func (s *Bolt) IterDB(fn func(key, value []byte) error) int64 {
	var total int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			if err := fn(key, value); err != nil {
				return err
			} else {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Bolt) IterKey(fn func(key []byte) error) int64 {
	var total int64
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		c := b.Cursor()
		for key, _ := c.First(); key != nil; key, _ = c.Next() {
			if err := fn(key); err != nil {
				return err
			} else {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (s *Bolt) Close() error {
	return s.db.Close()
}
