package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/vickon16/go-jwt-mysql/cmd/api"
	"github.com/vickon16/go-jwt-mysql/cmd/config"
	"github.com/vickon16/go-jwt-mysql/cmd/db"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DbUser,
		Passwd:               config.Envs.DbPassword,
		Addr:                 config.Envs.DbAddress,
		DBName:               config.Envs.DbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":"+config.Envs.Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
