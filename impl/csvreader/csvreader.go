/*Package csvreader provides an implementation of consumer.Reader interface for CSV as a source and some helpers*/
package csvreader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/alexbyk/goquiz/model"
)

// ParseCustomer convers array of strings to the Customer structure
func ParseCustomer(rec []string) (*model.Customer, error) {
	if len(rec) != 5 {
		return nil, errors.New("Wrong csv record")
	}
	return &model.Customer{ID: rec[0], FirstName: rec[1], LastName: rec[2], Email: rec[3], Phone: rec[4]}, nil
}

// ErrEmpty indicates EOF
var ErrEmpty = fmt.Errorf("Emty queue")

var validHeaderRecord = []string{"id", "first_name", "last_name", "email", "phone"}
var errInvalidHeader = fmt.Errorf("InvalidHeader")

// CheckCsvHeader returns an error if header isn't correct csv header
func CheckCsvHeader(rec []string) error {
	if len(rec) != len(validHeaderRecord) {
		return errInvalidHeader
	}
	for i, v := range validHeaderRecord {
		if rec[i] != v {
			return errInvalidHeader
		}
	}
	return nil
}

// Reader reads customers from the source which is csv input
type Reader struct {

	// How many records to read at most for each ReadRecords invocation
	readCount int

	// How many lines we already read
	headerChecked bool

	csv *csv.Reader
}

// NewReader creates and returns a Reader
func NewReader(r io.Reader, count int) *Reader {
	return &Reader{
		csv:       csv.NewReader(r),
		readCount: count,
	}
}

func (r *Reader) checkHeader() error {
	rec, err := r.csv.Read()
	if err != nil {
		return err
	}
	if err := CheckCsvHeader(rec); err != nil {
		return err
	}
	return nil
}

// ReadRecords returns a slice of records. Error may be io.EOF as well
// It reads as many correct records as posible. May return not empty array with error
func (r *Reader) ReadRecords() ([]*model.Customer, error) {

	var customers []*model.Customer
	var err error

	if !r.headerChecked {
		err = r.checkHeader()
		if err != nil {
			return customers, err
		}
		r.headerChecked = true
	}

	for count := 0; count < r.readCount; count++ {
		var rec []string
		rec, err = r.csv.Read()
		if err != nil {
			break
		}
		var customer *model.Customer
		customer, err = ParseCustomer(rec)
		if err != nil {
			break
		}
		customers = append(customers, customer)
	}
	return customers, err
}
