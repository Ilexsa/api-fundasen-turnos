package db

import(
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB


func ConnectDB(host string, port string, user string, pass string, dbName string, encrypt string) error {
	server := host
	if port != "" {
		server = fmt.Sprintf("%s,%s", host, port)
	}
	conString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=%s", server,
		user, pass, dbName, encrypt)
	var err error
	db, err = sql.Open("sqlserver", conString)
	if err!= nil{
		return  err
	}
	return db.Ping()
}

func GetDb() *sql.DB{
	return db
}