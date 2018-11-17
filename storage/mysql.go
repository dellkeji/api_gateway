package storage

import (
	"fmt"

	config "apigw_golang/configure"

	"github.com/jinzhu/gorm"

	// import _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBSession :
type DBSession struct {
	DB *gorm.DB
}

// Connect : connect the database
func (db *DBSession) Connect() error {
	var (
		baseConf = config.GlobalConfigurations
		dbConfig = baseConf.DBConf
		err      error
	)
	dbhost := fmt.Sprintf("tcp(%s:%d)", dbConfig.Host, dbConfig.Port)
	db.DB, err = gorm.Open(dbConfig.Type, fmt.Sprintf(
		"%s:%s@%s/%s?charset=%s&parseTime=True&loc=%s",
		dbConfig.User,
		dbConfig.Password,
		dbhost,
		dbConfig.DBName,
		dbConfig.Charset,
		baseConf.Location,
	))
	if err != nil {
		return err
	}

	sqldb := db.DB.DB()
	sqldb.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqldb.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqldb.Ping()

	if baseConf.Debug {
		db.DB.LogMode(true)
	}

	return nil

}

// Close :
func (db *DBSession) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
