package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"./datastores"
)

const (
	usrCookie = "UserID"
)

func main() {
	log.Println("Dental Web App: started.")
	controller := NewController()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/api/auth", controller.AuthHandler)
	http.HandleFunc("/api/appts", controller.GetApptsHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("FATAL ERROR: Could not listen for client requests; Details: %s", err)
	}
}

type Controller struct {
	userstore datastores.UserStore
	apptstore datastores.ApptStore
}

func NewController() *Controller {
	this := &Controller{
		userstore: datastores.NewMockUserDatastore(),
		apptstore: datastores.NewMockApptDatastore(),
	}

	//add default user
	user := &datastores.User{
		ID:        0,
		FirstName: "Admin",
		LastName:  "Jones",
		Email:     "test@test.com",
		UsrRole:   datastores.RoleAdmin,
	}
	this.userstore.AddUser(user)

	now := time.Now()

	appt := &datastores.Appt{
		ID:         0,
		Start:      now,
		End:        now.Add(time.Minute * 120),
		CustID:     1,
		HygID:      2,
		DenID:      3,
		CustName:   "Joe",
		HygName:    "Bob",
		DentName:   "Jill",
		ReqDentist: true,
	}
	this.apptstore.AddAppt(appt)

	return this
}

func (this *Controller) AuthHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		wtr.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(wtr, "%s: requst method was %q when POST was expected", http.StatusText(http.StatusMethodNotAllowed), req.Method)
		return
	}
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		wtr.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(wtr, "%s: request content type was not JSON", http.StatusText(http.StatusUnsupportedMediaType))
		return
	}
	authFields := make(map[string]string)

	if err := json.NewDecoder(req.Body).Decode(&authFields); err != nil {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: JSON unmarshall failed", http.StatusText(http.StatusBadRequest))
		return
	}
	defer req.Body.Close()

	email, passwd := authFields["Email"], authFields["Password"]
	fmt.Println(authFields)
	if email == "" || passwd == "" {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: username or password missing", http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := this.userstore.Auth(email, passwd)
	if err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: datastore error", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: Userdatastore failed to evaluate user auth request; Details: %s", err)
		return
	}
	if user == nil {
		wtr.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(wtr, "%s: access denied", http.StatusText(http.StatusForbidden))
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(wtr).Encode(user); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: JSON marshall failed", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: JSON marshall failed for user object; Details: %s", err)
		return
	}
	http.SetCookie(wtr, &http.Cookie{Name: usrCookie, Value: strconv.Itoa(user.ID), HttpOnly: true})
}

func (this *Controller) GetApptsHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		wtr.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(wtr, "%s: requst method was %q when GET was expected", http.StatusText(http.StatusMethodNotAllowed), req.Method)
		return
	}

	//TODO We will need this to filter later...
	/*usrCookieVal, err := req.Cookie(usrCookie)
	if err != nil {
		wtr.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(wtr, "%s: access denied, cookie missing", http.StatusText(http.StatusForbidden))
		return
	}

	/*usrID, err := strconv.Atoi(usrCookieVal.Value)
	if err != nil {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: user cookie invalid", http.StatusText(http.StatusBadRequest))
		return
	}*/

	params := req.URL.Query()

	startMs, err := strconv.ParseInt(params.Get("startDate"), 10, 64)
	if err != nil || startMs == 0 {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: start time is invalid", http.StatusText(http.StatusBadRequest))
		return
	}

	endMs, err := strconv.ParseInt(params.Get("endDate"), 10, 64)
	if err != nil || endMs == 0 {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: end time is invalid", http.StatusText(http.StatusBadRequest))
		return
	}

	appts, err := this.apptstore.GetApptByDate(time.Unix(0, startMs*1000000), time.Unix(0, endMs*1000000))
	if err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: datastore error", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: Apptdatastore failed to get appts; Details: %s", err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(wtr).Encode(appts); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: JSON marshall failed", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: JSON marshall failed for appts slice; Details: %s", err)
		return
	}
}

func (this *Controller) GetUsers(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		wtr.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(wtr, "%s: requst method was %q when GET was expected", http.StatusText(http.StatusMethodNotAllowed), req.Method)
		return
	}

	users, err := this.userstore.GetUsers()
	if err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: datastore error", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: Userdatastore failed to get users; Details: %s", err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(wtr).Encode(users); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: JSON marshall failed", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: JSON marshall failed for users slice; Details: %s", err)
		return
	}
}

func (this *Controller) AddApptHandler(wtr http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		wtr.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(wtr, "%s: requst method was %q when POST was expected", http.StatusText(http.StatusMethodNotAllowed), req.Method)
		return
	}

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		wtr.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(wtr, "%s: request content type was not JSON", http.StatusText(http.StatusUnsupportedMediaType))
		return
	}

	appt := new(datastores.Appt)
	err := json.NewDecoder(req.Body).Decode(&appt)
	if err != nil {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: JSON unmarshall failed", http.StatusText(http.StatusBadRequest))
		return
	}
	defer req.Body.Close()

	if err = this.apptstore.AddAppt(appt); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: datastore error", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: Apptdatastore failed to add appt; Details: %s", err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(wtr).Encode(appt); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: JSON marshall failed", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: JSON marshall failed for a new appt; Details: %s", err)
		return
	}
}

func (this *Controller) UpdateApptHandler(wtr http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		wtr.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(wtr, "%s: requst method was %q when POST was expected", http.StatusText(http.StatusMethodNotAllowed), req.Method)
		return
	}

	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		wtr.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(wtr, "%s: request content type was not JSON", http.StatusText(http.StatusUnsupportedMediaType))
		return
	}

	appt := new(datastores.Appt)
	err := json.NewDecoder(req.Body).Decode(&appt)
	if err != nil {
		wtr.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(wtr, "%s: JSON unmarshall failed", http.StatusText(http.StatusBadRequest))
		return
	}
	defer req.Body.Close()

	if err = this.apptstore.UpdateAppt(appt); err != nil {
		wtr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(wtr, "%s: datastore error", http.StatusText(http.StatusInternalServerError))
		log.Printf("ERROR: Apptdatastore failed to update appt; Details: %s", err)
		return
	}

	wtr.WriteHeader(http.StatusOK)
}
