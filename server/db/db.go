package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

func Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *gorm.DB {
	if Dbdriver != "postgres" {
		panic("only support postgres")
	}

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to the database %s", Dbdriver))
	}
	return db
}
