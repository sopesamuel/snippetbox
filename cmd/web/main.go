package main

import (

	"database/sql"
	"html/template"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.project.sope/internal/models"
	_ "github.com/go-sql-driver/mysql"

)

//Handlers
//Servermux
//Web server

//Holds depencies- sstuff our appliacation needs to work
type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}


func main(){
	
	//Configuration of terminal with flags
	addr := flag.String("addr", ":4000", "HTTP requests")
	dsn := flag.String("dsn", "web:pass@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name") 
	//looks for cli flags to change default
	flag.Parse() 

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	defer db.Close()
	logger.Info("Server starting", slog.String("addr", ":4000"))

	templateCache , err := newTemplateCache()
	if err != nil{
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger : logger,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string)(*sql.DB, error){

	db , err := sql.Open("mysql", dsn)
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