package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/DataDavD/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Additional info flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mysqlPw := os.Getenv("SNIPPETBOX_MYSQL_PW")
	pw := fmt.Sprintf("web:%s@/snippetbox?parseTime=true", mysqlPw)

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", pw, "MySQL data source name")
	flag.Parse()

	// To keep the main() func tidy we've put the code for creating a connection pool into separate
	// openDB() function below. We pass openDB() to the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			errorLog.Fatal(err)
		}
	}()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(), // Call the new app.routes() method
	}

	infoLog.Printf("Starting server on %s", *addr)
	if err = srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
