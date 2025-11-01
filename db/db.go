package db

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	 _ "github.com/go-sql-driver/mysql"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *sql.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	name := viper.GetString("database.name")
	port := viper.GetString("database.port")
	

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, name)

	db, err := sql.Open("mysql", dsn)
	if  err != nil {
		log.Fatalf("Failed to Initialize Database, Error: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database not reachable: %v", err)
	}

	return db
}
