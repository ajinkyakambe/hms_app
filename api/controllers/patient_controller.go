package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/repoerna/hms_app/api/auth"
	"github.com/repoerna/hms_app/api/handlers"
	"github.com/repoerna/hms_app/api/models"
	"github.com/repoerna/hms_app/utils/formaterror"
)

// CreatePatient ...
func (server *Server) CreatePatient(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
	}
	patient := models.Patient{}
	err = json.Unmarshal(body, &patient)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = patient.Validate("")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	patientCreated, err := patient.SavePatient(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, patientCreated.ID))
	handlers.ResponseJSON(w, http.StatusCreated, patientCreated)
}

// GetPatients ...
func (server *Server) GetPatients(w http.ResponseWriter, r *http.Request) {

	patient := models.Patient{}

	patients, err := patient.FindAllPatients(server.DB)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, patients)
}

// GetPatient ...
func (server *Server) GetPatient(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ssn, err := strconv.ParseUint(vars["ssn"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	patient := models.Patient{}
	patientGotten, err := patient.FindPatientBySSN(server.DB, uint32(ssn))
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, patientGotten)
}

// UpdatePatient ...
func (server *Server) UpdatePatient(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ssn, err := strconv.ParseUint(vars["ssn"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	patient := models.Patient{}
	err = json.Unmarshal(body, &patient)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(ssn) {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = patient.Validate("update")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := patient.UpdatePatient(server.DB, uint32(ssn))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, updatedUser)
}

// DeletePatient ...
func (server *Server) DeletePatient(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	patient := models.Patient{}

	ssn, err := strconv.ParseUint(vars["ssn"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(ssn) {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = patient.DeletePatient(server.DB, uint32(ssn))
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", ssn))
	handlers.ResponseJSON(w, http.StatusNoContent, "")
}
