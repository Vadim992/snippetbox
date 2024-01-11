package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/swaggest/form/v5"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.vadimpush.net/internal/models"
	"time"
)

type application struct {
	errLog         *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := openDB(*dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errLog:         errLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Starting server on port %s\n", *addr)

	errLog.Fatal(srv.ListenAndServe())

}

func openDB(dns string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dns)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
