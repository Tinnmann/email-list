package main

import (
	"database/sql"
	"email-list/edb"
	"email-list/grpcapi"
	"email-list/jsonapi"
	"log"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:EMAILLIST_DB"`
	BindJson string `arg:"env:EMAILLIST_BIND_JSON"`
	BindGrpc string `arg:"env:EMAILLIST_BIND_GRPC"`
}

func main() {
	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}
	if args.BindGrpc == "" {
		args.BindGrpc = ":8081"
	}

	log.Printf("using database '%v'\n", args.DbPath)
	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	edb.CreateDb(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting JSON API server...\n")
		jsonapi.Serve(db, args.BindJson)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Printf("starting gRPC API server...\n")
		grpcapi.Serve(db, args.BindGrpc)
		wg.Done()
	}()

	wg.Wait()
}
