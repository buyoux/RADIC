package kvdb

import (
	"RADIC/util"
	"os"
	"strings"
)

// 几种常见的基于LSM-Tree算法实现的KV数据库
const (
	BOLT = iota
	BADGER
)

// redis也是一种KV数据库，读者可以自行用redis实现IKeyValueDB接口
type IKeyValueDB interface {
	Open() error                                   // 初始化DB
	GetDbPath() string                             // 获取存储数据库的目录
	Set(key, value []byte) error                   // 写入<key, value>, k,v为任意长度的字节流
	BatchSet(keys, values [][]byte) error          // 批量写入<key, value>
	Get(key []byte) ([]byte, error)                // 读取key对应的value
	BatchGet(keys [][]byte) ([][]byte, error)      // 批量读取key对应的value
	Delete(key []byte) error                       // 根据key删除对应kv
	BatchDelete(keys [][]byte) error               // 批量删除
	Has(key []byte) bool                           // 判断某个key是否存在
	IterDB(fn func(key, value []byte) error) int64 // 遍历数据库，返回数据条数
	IterKey(fn func(key []byte) error) int64       // 遍历所有key，返回数据条数
	Close() error                                  // 把内存中的数据flush到磁盘，同时释放文件锁
}

// Factory工厂模式，把类的创建和使用分隔开。
// Get函数就是一个工厂，它返回产品的接口，即它可以返回各种各样的具体产品（数据库）
func GetKvDb(dbtype int, path string) (IKeyValueDB, error) {
	paths := strings.Split(path, "/")
	parentPath := strings.Join(paths[0:len(paths)-1], "/") // 父路径

	info, err := os.Stat(parentPath)
	// 若父路径不存在则创建路径
	if os.IsNotExist(err) {
		util.Log.Printf("create dir %s", parentPath)
		// 数字前的0或0o都表示八进制 600代表权限
		os.MkdirAll(parentPath, 0o600)
	} else { // 父路径存在
		// 若父路径为一个普通文件，则移除
		if info.Mode().IsRegular() {
			util.Log.Printf("%s is a regular file, will delete it", parentPath)
			os.Remove(parentPath)
		}
	}

	var db IKeyValueDB
	switch dbtype {
	case BADGER:
		db = new(Badger).WithDataPath(path)
	default: // 默认使用bolt
		// Builder生成器模式
		db = new(Bolt).WithDataPath(path).WithBucket("radic")
	}
	// 创建具体KVDB的细节隐藏在Open()函数中，在这里创建类
	err = db.Open()
	return db, err
}
