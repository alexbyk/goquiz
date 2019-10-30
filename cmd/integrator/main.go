package main

import (
	"fmt"
	"os"

	"github.com/alexbyk/goquiz/common/integrator"
	"github.com/alexbyk/goquiz/impl/httpapi"
	"github.com/alexbyk/goquiz/impl/pgstorage"
)

var defaultEndpoint = "http://localhost:8080"

const defaultDSN = "postgres://localhost/postgres?sslmode=disable"

func fail(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}

func checkError(err error) {
	if err != nil {
		fail(err)
	}
}

func main() {
	var endpoint, dsn string
	if endpoint = os.Getenv("API_ENDPOINT"); endpoint == "" {
		endpoint = defaultEndpoint
	}

	if dsn = os.Getenv("DSN"); dsn == "" {
		dsn = defaultDSN
	}

	db, err := pgstorage.Connect(dsn)
	checkError(err)
	checkError(pgstorage.CreateTable(db))
	defer db.Close()

	integrator := integrator.NewIntegrator(pgstorage.NewPgNexter(db), pgstorage.NewPgListener(db), httpapi.NewHTTPPublisher(endpoint))
	integrator.Start()
}
