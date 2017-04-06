package datastores

import (
	"time"
)

type ApptStore interface {
	AddAppt(*Appt) error
	DelAppt(ID int) error
	UpdateAppt(*Appt) error
	GetApptByDate(start, end time.Time) ([]Appt, error)
}

type UserStore interface {
	AddUser(*User) error
	DelUser(ID int) error
	UpdateUser(*User) error
	Auth(email, password string) (*User, error)
	GetUsers() ([]User, error)
}
