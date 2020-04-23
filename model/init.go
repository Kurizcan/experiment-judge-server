package model

import (
	"fmt"

	"github.com/lexkong/log"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Self        *gorm.DB
	ApiDatabase *gorm.DB
}

var DB *Database

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s&multiStatements=true",
		username,
		password,
		addr,
		name,
		true,
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	}

	// set for db connection
	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

// used for cli
func InitSelfDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

func InitApiDB() *gorm.DB {
	return openDB(viper.GetString("api_server_db.username"),
		viper.GetString("api_server_db.password"),
		viper.GetString("api_server_db.addr"),
		viper.GetString("api_server_db.name"))
}

func GetApiDB() *gorm.DB {
	return InitApiDB()
}

func (db *Database) Init() {
	DB = &Database{
		Self:        GetSelfDB(),
		ApiDatabase: GetApiDB(),
	}
}

func (db *Database) Close() {
	DB.Self.Close()
	DB.ApiDatabase.Close()
}

func CreateDataBase(dataBase string) error {
	sql := fmt.Sprintf("create database if not exists %s", dataBase)
	return DB.Self.Exec(sql).Error
}
