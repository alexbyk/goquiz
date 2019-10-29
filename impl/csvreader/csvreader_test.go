package csvreader_test

import (
	"io"
	"strings"
	"testing"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/goquiz/impl/csvreader"
	"github.com/alexbyk/goquiz/common/model"
)

func TestParseCustomer(t *testing.T) {
	ft := ftest.New(t)
	okArr := []string{"2", "alex", "byk", "alex@alexbyk.com", "38066 99 18018"}
	okRec := &model.Customer{ID: "2", FirstName: "alex", LastName: "byk", Email: "alex@alexbyk.com", Phone: "38066 99 18018"}

	// ok
	rec, err := csvreader.ParseCustomer(okArr)
	ft.Nil(err).Eq(rec, okRec)

	// fail
	rec, err = csvreader.ParseCustomer(okArr[1:])
	ft.Nil(rec).NotNil(err)

	rec, err = csvreader.ParseCustomer(append(okArr, "foo"))
	ft.Nil(rec).NotNil(err)
}

func TestIsValidCsvHeader(t *testing.T) {
	ft := ftest.New(t)
	ok := []string{"id", "first_name", "last_name", "email", "phone"}
	ft.Nil(csvreader.CheckCsvHeader(ok))
	ft.NotNil(csvreader.CheckCsvHeader(nil))
	ft.NotNil(csvreader.CheckCsvHeader(append(ok, "extra")))
	ft.NotNil(csvreader.CheckCsvHeader(ok[1:]))
}

func TestReadRecords(t *testing.T) {

	ft := ftest.New(t)
	in := `id,first_name,last_name,email,phone
1,alex,byk,alex@alexbyk.com,38066 99 18018
2,alex,byk,alex@alexbyk.com,38066 99 18018
3,alex,byk,alex@alexbyk.com,38066 99 18018
`
	tests := []struct {
		customers []*model.Customer
		err       error
	}{
		{[]*model.Customer{{ID: "1", FirstName: "alex", LastName: "byk", Email: "alex@alexbyk.com", Phone: "38066 99 18018"}, {ID: "2", FirstName: "alex", LastName: "byk", Email: "alex@alexbyk.com", Phone: "38066 99 18018"}}, nil},
		{[]*model.Customer{{ID: "3", FirstName: "alex", LastName: "byk", Email: "alex@alexbyk.com", Phone: "38066 99 18018"}}, io.EOF},
	}

	r := csvreader.NewReader(strings.NewReader(in), 2)

	for _, test := range tests {
		records, err := r.ReadRecords()
		ft.Eq(err, test.err).Eq(records, test.customers)
	}

}

func TestReadRecordFail(t *testing.T) {
	ft := ftest.New(t)
	in := `id,first_name,last_name,email
1,alex,byk,alex@alexbyk.com
`
	r := csvreader.NewReader(strings.NewReader(in), 2)
	_, err := r.ReadRecords()
	ft.NotNil(err)
}
