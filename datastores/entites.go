package datastores

import (
	"time"
)

const (
	RoleAdmin = iota
	RoleCust
	RoleHyg
	RoleDen
)

type Role uint8

type Appt struct {
	ID                          int
	Start, End                  time.Time
	ReqDentist                  bool
	CustID, HygID, DenID        int
	CustName, HygName, DentName string
}

type User struct {
	ID                         int
	LastName, FirstName, Email string
	UsrRole                    Role
}
