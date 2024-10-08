// Package mysql_config provides functionality for configuring MySQL database connection.
package mysql_config

import (
	"database/sql"
	"fmt"
	"os"
	awsconfig "test/lambda/aws-config"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabase represents a MySQL database connection.
type MysqlDatabase struct {
	Db *sql.DB
}

// GetConn returns the underlying SQL database connection.
func (mysql *MysqlDatabase) GetConn() *sql.DB {
	return mysql.Db
}

// NewMysqlDatabase creates a new MySQL database instance.
// It fetches database credentials from AWS Secrets Manager based on the provided secret ID.
// It returns a pointer to the mysqlDatabase or panics if an error occurs during initialization.
func NewMysqlDatabase() *MysqlDatabase {
	env, err := awsconfig.GetSecret(os.Getenv("SECRET_ID_MYSQL"))
	if err != nil {
		fmt.Println("Error fetching secret", err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		env["DB_USERNAME_MYSQL"],
		env["DB_PASSWORD_MYSQL"],
		env["DB_HOST_MYSQL"],
		env["DB_PORT_MYSQL"],
		env["DB_DATABASE_MYSQL"],
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return &MysqlDatabase{Db: db}
}
