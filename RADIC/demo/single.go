package course

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 单例模式

var single *gorm.DB // 通过gorm.Open()创建的gorm.DB是一个连接池，只需创建一次实例
var once sync.Once = sync.Once{}
var lock = &sync.Mutex{}

func GetDB1() *gorm.DB {
	if single == nil {
		lock.Lock()
		defer lock.Unlock()
		if single == nil {
			single, _ = gorm.Open(mysql.Open(""))
		} else {
			fmt.Println("数据库连接池单例已经创建")
		}
	} else {
		fmt.Println("数据库连接池单例已经创建")
	}
	return single
}

// init()只会执行一次，可以用来实现单例。
// 但使用init()通常要小心代码的依赖关系
func init() {
	single, _ = gorm.Open(mysql.Open(""))
}

func GetDB2() *gorm.DB {
	return single
}

func GetDB3() *gorm.DB {
	if single == nil {
		once.Do(func() {
			single, _ = gorm.Open(mysql.Open(""))
		})
	} else {
		fmt.Println("数据库连接池单例已经创建")
	}
	return single
}
