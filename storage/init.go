package storage

// once synchronization
import (
	"sync"
)

// StorageSession : interface of db session
type StorageSession interface {
	Connect() error
	Close()
}

// GLobals
var (
	GlobalDBSession *DBSession
)

var dbOnce sync.Once

// GetDBSession : init the db session pool
func GetDBSession() *DBSession {
	if GlobalDBSession == nil {
		dbOnce.Do(func() {
			GlobalDBSession = &DBSession{}
			err := GlobalDBSession.Connect()
			if err != nil {
				panic(err)
			}
		})
	}
	return GlobalDBSession
}
