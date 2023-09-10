package monitor

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitDB() {
	const (
		host     = "localhost"
		port     = 3306
		user     = "root"
		password = "123456"
		dbname   = "serverless"
	)
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", mysqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func LoggerPerformance(port int, startTime string, spend float64, requestParams string, url string, responseBodyLength, statusCode int) {

	_, err := DB.Exec(`INSERT INTO 
    performance_monitoring 
    (port, request_start_time, response_time_ms, request_parameters, response_body_length, response_status_code, request_url)
	VALUES 
	    (?, ?, ?, ?, ?, ?, ?)`,
		port, startTime, spend, requestParams, responseBodyLength, statusCode, url)

	if err != nil {
		log.Printf("Failed to log performance data: %v", err)
	}
}
