package database

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gocraft/dbr/v2"
)

var Sqlinstance = map[string]*dbr.Connection{}
var SqlsessionError = map[string]error{}
var SqlonceMap = map[string]*sync.Once{}
var lock sync.RWMutex

func GetSessionForSQL(SQLDsn string) (*dbr.Session, error) {
	lock.RLock()
	once, ok := SqlonceMap[SQLDsn]
	lock.RUnlock()
	if !ok {
		once = &sync.Once{}
		once.Do(func() {
			lock.Lock()
			defer lock.Unlock()
			SqlonceMap[SQLDsn] = once
			connection, err := dbr.Open("mysql", SQLDsn, nil)
			if err != nil {
				SqlsessionError[SQLDsn] = err
				log.Println("Error - getSQLConnection - ", err)
			}
			Sqlinstance[SQLDsn] = connection

			poolSetting, err := strconv.Atoi(os.Getenv("maxidleconns"))
			if err != nil {
				log.Println("Error - maxidleconns: ", err)
			}
			Sqlinstance[SQLDsn].SetMaxIdleConns(poolSetting)
			//Get max open connection from config
			poolSetting, err = strconv.Atoi(os.Getenv("maxopenconns"))
			if err != nil {
				log.Println("Error - maxopenconns: ", err)
			}
			Sqlinstance[SQLDsn].SetMaxOpenConns(poolSetting)
		})
	}
	log.Println("connected to SQLConnetcion!!!")
	//Get max idle connection from config
	return Sqlinstance[SQLDsn].NewSession(nil), SqlsessionError[SQLDsn]
}
