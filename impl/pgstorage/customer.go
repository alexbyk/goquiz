package pgstorage

import (
	"github.com/alexbyk/goquiz/common/model"
)

// Customer inherits model.Customer because we need to maintain a status of the record and not depend on the implementation
type Customer struct {
	*model.Customer
	Status string `sql:",notnull"`
}

const (
	//CustomerStatusEmpty customer that should be send
	CustomerStatusEmpty = ""

	//CustomerStatusConfirmed means succesfully sent
	CustomerStatusConfirmed = "confirmed"

	//CustomerStatusSending means we're processing it right now, and if there are records with such status, it should be handled separately
	CustomerStatusSending = "sending"
)
