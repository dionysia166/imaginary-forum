package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"forum/cmd/internal/models"

	"github.com/alexedwards/scs/v2"   
	_ "github.com/mattn/go-sqlite3"
)

// application holds the application-wide dependencies.
type application struct {
	logger        *slog.Logger
	messages      *models.MessageModel
	threads       *models.ThreadModel
	users         *models.UserModel
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dbPath := flag.String("db", "./db.sqlite", "Path to SQLite database")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dbPath)
	if err != nil {
		logger.Error((err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:        logger,
		messages:      &models.MessageModel{DB: db},
		threads:       &models.ThreadModel{DB: db},
		users:         &models.UserModel{DB: db},
		templateCache: templateCache,
		sessionManager: sessionManager,
	}

	logger.Info("Starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB opens a connection to the SQLite database.
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
