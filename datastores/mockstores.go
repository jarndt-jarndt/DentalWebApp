package datastores

import (
	"errors"
	"time"
)

var _ ApptStore = new(MockApptStore)

type MockApptStore struct {
	nextId int
	appts  []Appt
}

func NewMockApptDatastore() *MockApptStore {
	return &MockApptStore{
		appts: make([]Appt, 0, 8),
	}
}

func (this *MockApptStore) AddAppt(appt *Appt) error {
	appt.ID = this.nextId
	this.nextId++
	this.appts = append(this.appts, *appt)
	return nil
}

func (this *MockApptStore) DelAppt(ID int) error {
	for i, appt := range this.appts {
		if appt.ID == ID {
			this.appts = append(this.appts[:i], this.appts[i+1:]...)
			return nil
		}
	}

	return nil
}

func (this *MockApptStore) UpdateAppt(updatedAppt *Appt) error {
	for i, appt := range this.appts {
		if appt.ID == updatedAppt.ID {
			this.appts[i] = *updatedAppt
			return nil
		}
	}
	return errors.New("Could not update appt, does not exsit")
}

func (this *MockApptStore) GetApptByDate(start, end time.Time) ([]Appt, error) {
	results := make([]Appt, 0, len(this.appts))
	for _, appt := range this.appts {
		if (appt.Start.After(start) || appt.Start.Equal(start)) &&
			(appt.End.Before(end) || appt.End.Equal(end)) {
			results = append(results, appt)
		}
	}
	return results, nil
}

var _ UserStore = new(MockUserStore)

type MockUserStore struct {
	users []User
}

func NewMockUserDatastore() *MockUserStore {
	return &MockUserStore{
		users: make([]User, 0, 8),
	}
}

func (this *MockUserStore) AddUser(user *User) error {
	this.users = append(this.users, *user)
	return nil
}

func (this *MockUserStore) DelUser(ID int) error {
	for i, user := range this.users {
		if user.ID == ID {
			this.users = append(this.users[:i], this.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func (this *MockUserStore) UpdateUser(updatedUser *User) error {
	for i, user := range this.users {
		if user.ID == updatedUser.ID {
			this.users[i] = *updatedUser
			return nil
		}
	}
	return errors.New("Could not update user, does not exist")
}

func (this *MockUserStore) Auth(email, password string) (*User, error) {
	for _, user := range this.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (this *MockUserStore) GetUsers() ([]User, error) {
	return this.users, nil
}
