package dao

import (
	"fmt"
	"sync"
	"time"

	"goframework/internal/dao"

	"github.com/go-sql-driver/mysql"
	"github.com/jianbo-zh/go-config"
	"github.com/jianbo-zh/go-errors"
	"gorm.io/gorm"
)

var gdb *gorm.DB

var lockgdb sync.Mutex

func DB() *gorm.DB {
	return MysqlDB()
}

func MysqlDB() *gorm.DB {
	if gdb == nil {

		lockgdb.Lock()
		defer lockgdb.Unlock()

		if gdb == nil {
			myconfig := mysql.Config{
				Net: "tcp",
				Addr: fmt.Sprintf(
					"%s:%s",
					config.GetString("mysql-main.host"),
					config.GetString("mysql-main.password"),
				),
				User:   config.GetString("mysql-main.username"),
				Passwd: config.GetString("mysql-main.password"),
				DBName: config.GetString("mysql-main.database"),
				Params: map[string]string{
					"charset": config.GetString("mysql-main.charset"),
				},
				ParseTime: true,
				Loc:       time.Local,
			}
			var err error
			gdb, err = dao.Mysql(myconfig.FormatDSN())
			if err != nil {
				panic(errors.Newf("get main db error: %s", err))
			}
		}
	}

	return gdb
}
