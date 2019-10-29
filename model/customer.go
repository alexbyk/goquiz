package model

// Customer represents a customer model
type Customer struct {
	ID        string `sql:",notnull"`
	FirstName string `sql:",notnull"`
	LastName  string `sql:",notnull"`
	Email     string `sql:",notnull"`
	Phone     string `sql:",notnull"`

	// Confirmed means that this record was successfully sent to crm
	Confirmed bool `sql:",notnull"`
}
