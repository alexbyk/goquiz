package main

import (
	"fmt"
	"os"

	"github.com/alexbyk/goquiz/common/consumer"

	"github.com/alexbyk/goquiz/impl/csvreader"
	"github.com/alexbyk/goquiz/impl/pgstorage"
)

const defaultDSN = "postgres://localhost/postgres?sslmode=disable"
const defaultFile = "data.csv"
const chunkSize = 10 // 10 lines per once
var usageMsg = fmt.Sprintf(`DSN=%s ./parser %s`, defaultDSN, defaultFile)

func fail(msg interface{}) {
	fmt.Println(msg, "\n", "Usage example:", "\n\t", usageMsg)
	os.Exit(1)
}

func checkError(err error) {
	if err != nil {
		fail(err)
	}
}

func main() {
	var dsn, filepath string
	if dsn = os.Getenv("DSN"); dsn == "" {
		dsn = defaultDSN
	}

	if len(os.Args) < 2 {
		filepath = defaultFile
	} else {
		filepath = os.Args[1]
	}

	db, err := pgstorage.Connect(dsn)
	checkError(err)
	checkError(pgstorage.CreateTable(db))
	defer db.Close()

	f, err := os.Open(filepath)
	checkError(err)
	defer f.Close()

	consumer := consumer.NewConsumer(csvreader.NewReader(f, chunkSize), pgstorage.NewPgWriter(db), pgstorage.NewPgNotifier(db))
	checkError(consumer.Consume())
}
