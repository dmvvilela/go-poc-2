package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dmvvilela/go-poc-2/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"name":"Pet", "email": "pet@gmail.com"}`,
			statusCode:   201,
			name:         "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Frank", "email": "pet@gmail.com"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"name":"Pet", "email": "grand@gmail.com"}`,
			statusCode:   201,
			name:         "Pet",
			email:        "grand@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"name":"Kan", "email": "kangmail.com"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"name": "", "email": "kan@gmail.com"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			inputJSON:    `{"name": "Kan", "email": ""}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/contacts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateContact)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetContacts(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedContacts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/contacts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetContacts)
	handler.ServeHTTP(rr, req)

	var contacts []models.Contact
	err = json.Unmarshal([]byte(rr.Body.String()), &contacts)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(contacts), 2)
}

func TestGetContactByID(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contact, err := seedOneContact()
	if err != nil {
		log.Fatal(err)
	}

	contactSample := []struct {
		id           string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(contact.ID)),
			statusCode: 200,
			name:       contact.Name,
			email:      contact.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range contactSample {
		req, err := http.NewRequest("GET", "/contacts", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetContact)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, contact.Name, responseMap["name"])
			assert.Equal(t, contact.Email, responseMap["email"])
		}
	}
}

func TestUpdateContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contacts, err := seedContacts() //we need atleast two contacts to properly check the update
	if err != nil {
		log.Fatalf("Error seeding contact: %v\n", err)
	}

	// Get only the first contact
	contact := contacts[0]
	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateName   string
		updateEmail  string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(contact.ID)),
			updateJSON:   `{"name":"Grand", "email": "grand@gmail.com"}`,
			statusCode:   200,
			updateName:   "Grand",
			updateEmail:  "grand@gmail.com",
			errorMessage: "",
		},
		{
			// Remember "kenny@gmail.com" belongs to contact 2
			id:           strconv.Itoa(int(contact.ID)),
			updateJSON:   `{"name":"Frank", "email": "kenny@gmail.com"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			id:           strconv.Itoa(int(contact.ID)),
			updateJSON:   `{"name":"Kan", "email": "kangmail.com"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(contact.ID)),
			updateJSON:   `{"name": "", "email": "kan@gmail.com"}`,
			statusCode:   422,
			errorMessage: "Required Name",
		},
		{
			id:           strconv.Itoa(int(contact.ID)),
			updateJSON:   `{"name": "Kan", "email": ""}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/contacts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateContact)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["name"], v.updateName)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contacts, err := seedContacts() //we need atleast two contacts to properly check the update
	if err != nil {
		log.Fatalf("Error seeding contact: %v\n", err)
	}

	// Get only the first contact
	contact := contacts[0]
	contactSample := []struct {
		id           string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(contact.ID)),
			statusCode:   204,
			errorMessage: "",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range contactSample {
		req, err := http.NewRequest("GET", "/contacts", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteContact)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
