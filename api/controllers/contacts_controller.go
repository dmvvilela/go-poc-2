package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dmvvilela/go-poc-2/api/models"
	"github.com/dmvvilela/go-poc-2/api/responses"
	"github.com/dmvvilela/go-poc-2/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateContact(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	contact := models.Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	contact.Prepare()
	err = contact.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	contactCreated, err := contact.SaveContact(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, contactCreated.ID))
	responses.JSON(w, http.StatusCreated, contactCreated)
}

func (server *Server) GetContacts(w http.ResponseWriter, r *http.Request) {
	contact := models.Contact{}

	contacts, err := contact.FindAllContacts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, contacts)
}

func (server *Server) GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	contact := models.Contact{}
	contactGotten, err := contact.FindContactByID(server.DB, uint64(cid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, contactGotten)
}

func (server *Server) UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	contact := models.Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	contact.Prepare()
	err = contact.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedContact, err := contact.UpdateContact(server.DB, uint32(cid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedContact)
}

func (server *Server) DeleteContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contact := models.Contact{}

	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = contact.DeleteContact(server.DB, uint64(cid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
