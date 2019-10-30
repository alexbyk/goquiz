package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/alexbyk/goquiz/common/consumer"

	"github.com/alexbyk/goquiz/impl/csvreader"
	"github.com/alexbyk/goquiz/impl/pgstorage"
)

const defaultDSN = "postgres://localhost/postgres?sslmode=disable"
const chunkSize = 10 // 10 lines per once
var usageMsg = fmt.Sprintf(`DSN=%s ./parser ./data.csv`, defaultDSN)

var mockFileSize = 10000

func fail(msg interface{}) {
	fmt.Println(msg, "\n", "Usage example:", "\n\t", usageMsg)
	os.Exit(1)
}

func checkError(err error) {
	if err != nil {
		fail(err)
	}
}

// For demo purposes if no file is provided as an argument, we generate a temporary csv file and use it
func main() {
	var dsn, filepath string
	if dsn = os.Getenv("DSN"); dsn == "" {
		dsn = defaultDSN
	}

	if len(os.Args) < 2 {
		file := generateMock(mockFileSize)
		defer os.Remove(file.Name())
		filepath = file.Name()
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

	cr := consumer.NewConsumer(csvreader.NewReader(f, chunkSize), pgstorage.NewPgWriter(db), pgstorage.NewPgNotifier(db))
	log.Printf("Parsing CSV file, %s", f.Name())
	checkError(cr.Consume())
	log.Println("Done")
}

var tmpl = []string{"alex", "byk", "alex@alexbyk.com", "38066 99 18018"}

func generateMock(n int) *os.File {
	tmpfile, err := ioutil.TempFile("", "example")
	log.Printf("Generating mock content in %s for %v records", tmpfile.Name(), n)
	if err != nil {
		fail(err)
	}

	wr := csv.NewWriter(tmpfile)
	wr.Write(csvreader.ValidHeaderRecord) // header

	for i := 0; i <= n; i++ {
		rec := append([]string{strconv.Itoa(i)}, tmpl...)
		e := wr.Write(rec)
		if e != nil {
			fail(e)
		}

	}
	wr.Flush()
	log.Printf("Finished %v records", n)
	return tmpfile
}
