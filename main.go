package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
	db      *sql.DB
}

func main() {

	server_database_dsn := "aryan_zcopy:xyab*!2j6imr@tcp(mysql-aryan.alwaysdata.net:3306)/aryan_zcopy"
	port := ":8080"

	// local_database_dsn := "zcpy:pass@/zcopy?parseTime=true"
	// port := ":8080"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("mysql", server_database_dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog: infoLog,
		errLog:  errLog,
		db:      db,
	}

	srv := &http.Server{
		Addr:     port,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Server started on port 8080")
	log.Fatal(srv.ListenAndServe())

}
