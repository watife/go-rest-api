package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// dynamic http address from command-line flag with default :4000
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse the command-line flag
	flag.Parse()

	// Info log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Error log
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Lshortfile) //can also tage log.Llongfile for full path

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// custom server struct to make use of custom errorLog
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// Listen to server
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)
}
