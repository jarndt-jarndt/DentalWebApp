package datastores

import (
	"testing"
	"time"
)

func TestMockApptStoreBasic(t *testing.T) {
	testApptStoreBasic(NewMockApptDatastore(), t)
}

func TestMockUserStoreBasic(t *testing.T) {
	testUserStoreBasic(NewMockUserDatastore(), t)
}

func testApptStoreBasic(store ApptStore, t *testing.T) {
	now := time.Now()
	appt := &Appt{
		ID:    0,
		Start: now,
		End:   now.Add(time.Minute * 120),
	}

	//trying to add created appt to datastore
	err := store.AddAppt(appt)
	if err != nil {
		t.Fatalf("Could not add appt to datastore; Details: %s", err)
	}

	//changing the DenID in created appt to 1 then trying to update the appt in datastore
	appt.DenID = 1
	if err = store.UpdateAppt(appt); err != nil {
		t.Fatalf("Could not update appt in datastore; Details: %s", err)
	}

	//trying to get an appt that is not in the range of exsiting appt
	results, err := store.GetApptByDate(now.Add(time.Minute*121), now.Add(time.Minute*123))
	if err != nil {
		t.Fatalf("Could not get appts from datastore; Details: %s", err)
	}
	if len(results) != 0 {
		t.Fatalf("%d appts found in unexcepted range!", len(results))
	}

	//trying to get an appt that is exactly the range of the exisiting appt
	results, err = store.GetApptByDate(now, now.Add(time.Minute*120))
	if err != nil {
		t.Fatalf("Could not get appts from datastore; Details: %s", err)
	}
	if len(results) != 1 {
		t.Fatalf("excepted 1 appt in range, found %d in range!", len(results))
	}

	//Deletes the exsiting appt by the appt ID
	if err = store.DelAppt(appt.ID); err != nil {
		t.Fatalf("Could not remove appt from datastore: Details: %s", err)
	}

	//trys to get the deleted appt
	results, err = store.GetApptByDate(now, now.Add(time.Minute*120))
	if err != nil {
		t.Fatalf("Could not get appts from datastore; Details: %s", err)
	}
	if len(results) != 0 {
		t.Fatalf("excepted 0 appt in range, found %d in range!", len(results))
	}
}

func testUserStoreBasic(store UserStore, t *testing.T) {
	user := &User{
		ID:        0,
		FirstName: "Ender",
		LastName:  "Jones",
		Email:     "test@test.com",
	}

	// Trys to add a user to the data store
	err := store.AddUser(user)
	if err != nil {
		t.Fatalf("Could not add user to datastore: Details: %s", err)
	}

	//Changes the user ID and then updates the the ID in the data store
	user.FirstName = "Bob"
	if err = store.UpdateUser(user); err != nil {
		t.Fatalf("Could not update user in datastore; Details: %s", err)
	}

	//Deletes current user
	if err := store.DelUser(user.ID); err != nil {
		t.Fatalf("Could not remove user from datastore: Details: %s", err)
	}

	//trys to update the previously deleted user
	user.Email = "foo@test.com"
	if err := store.UpdateUser(user); err == nil {
		t.Fatalf("User remained in data store after deleting: Details: %s", err)
	}
}
