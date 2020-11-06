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
	"github.com/repoerna/hms_app/api/utils/formaterror"
)

// CreateExamination ...
func (server *Server) CreateExamination(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
	}
	examination := models.Examination{}
	err = json.Unmarshal(body, &examination)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = examination.Validate("")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	examinationCreated, err := examination.SaveExamination(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, examinationCreated.ExaminationID))
	handlers.ResponseJSON(w, http.StatusCreated, examinationCreated)
}

// GetExaminations ...
func (server *Server) GetExaminations(w http.ResponseWriter, r *http.Request) {

	examination := models.Examination{}

	examinations, err := examination.FindAllExamination(server.DB)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, examinations)
}

// GetExamination ...
func (server *Server) GetExamination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eid, err := strconv.ParseUint(vars["examination_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	examination := models.Examination{}
	examinationGotten, err := examination.FindExaminationByID(server.DB, uint32(eid))
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, examinationGotten)
}

// UpdateExamination ...
func (server *Server) UpdateExamination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)

	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	eid, err := strconv.ParseUint(vars["examination_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	examination := models.Examination{}
	err = json.Unmarshal(body, &examination)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = examination.Validate("update")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedExamination, err := examination.UpdateExamination(server.DB, uint32(eid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, updatedExamination)
}

// DeleteExamination ...
func (server *Server) DeleteExamination(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	examination := models.Examination{}

	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	eid, err := strconv.ParseUint(vars["examination_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = examination.DeleteExamination(server.DB, uint32(uid))
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", eid))
	handlers.ResponseJSON(w, http.StatusNoContent, "")
}
