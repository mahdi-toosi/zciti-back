package schema

import (
	"github.com/lib/pq"
)

type UserView struct {
	ID              uint64
	FirstName       string
	LastName        string
	FullName        string
	Mobile          uint64
	MobileConfirmed bool
	Roles           pq.StringArray
	Base
}
