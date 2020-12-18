package modeltests

import (
	"log"
	"testing"

	"github.com/dmvvilela/go-poc-2/api/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllContacts(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedContacts()
	if err != nil {
		log.Fatal(err)
	}

	contacts, err := contactInstance.FindAllContacts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the contacts: %v\n", err)
		return
	}

	assert.Equal(t, len(*contacts), 2)
}

func TestSaveContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	newContact := models.Contact{
		ID:    1,
		Name:  "test",
		Email: "test@gmail.com",
	}

	savedContact, err := newContact.SaveContact(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the contacts: %v\n", err)
		return
	}

	assert.Equal(t, newContact.ID, savedContact.ID)
	assert.Equal(t, newContact.Email, savedContact.Email)
	assert.Equal(t, newContact.Name, savedContact.Name)
}

func TestGetContactByID(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contact, err := seedOneContact()
	if err != nil {
		log.Fatalf("cannot seed contacts table: %v", err)
	}

	foundContact, err := contactInstance.FindContactByID(server.DB, contact.ID)
	if err != nil {
		t.Errorf("this is the error getting one contact: %v\n", err)
		return
	}

	assert.Equal(t, foundContact.ID, contact.ID)
	assert.Equal(t, foundContact.Email, contact.Email)
	assert.Equal(t, foundContact.Name, contact.Name)
}

func TestUpdateContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contact, err := seedOneContact()
	if err != nil {
		log.Fatalf("Cannot seed contact: %v\n", err)
	}

	contactUpdate := models.Contact{
		ID:    1,
		Name:  "modiUpdate",
		Email: "modiupdate@gmail.com",
	}

	updatedContact, err := contactUpdate.UpdateContact(server.DB, uint32(contact.ID))
	if err != nil {
		t.Errorf("this is the error updating the contact: %v\n", err)
		return
	}

	assert.Equal(t, updatedContact.ID, contactUpdate.ID)
	assert.Equal(t, updatedContact.Email, contactUpdate.Email)
	assert.Equal(t, updatedContact.Name, contactUpdate.Name)
}

func TestDeleteContact(t *testing.T) {
	err := refreshContactTable()
	if err != nil {
		log.Fatal(err)
	}

	contact, err := seedOneContact()
	if err != nil {
		log.Fatalf("Cannot seed contact: %v\n", err)
	}

	isDeleted, err := contactInstance.DeleteContact(server.DB, contact.ID)
	if err != nil {
		t.Errorf("this is the error updating the contact: %v\n", err)
		return
	}

	// one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	// Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
