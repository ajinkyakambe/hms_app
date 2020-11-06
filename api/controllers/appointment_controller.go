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

// CreateAppointment ...
func (server *Server) CreateAppointment(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
	}
	appointment := models.Appointment{}
	err = json.Unmarshal(body, &appointment)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = appointment.Validate("")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	appointmentCreated, err := appointment.SaveAppointment(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, appointmentCreated.AppointmentID))
	handlers.ResponseJSON(w, http.StatusCreated, appointmentCreated)
}

// GetAppointments ...
func (server *Server) GetAppointments(w http.ResponseWriter, r *http.Request) {

	appointment := models.Appointment{}

	appointments, err := appointment.FindAllAppointment(server.DB)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, appointments)
}

// GetAppointment ...
func (server *Server) GetAppointment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	aid, err := strconv.ParseUint(vars["appointment_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	appointment := models.Appointment{}
	appointmentGotten, err := appointment.FindAppointmentByID(server.DB, uint32(aid))
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, appointmentGotten)
}

// UpdateAppointment ...
func (server *Server) UpdateAppointment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)

	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	aid, err := strconv.ParseUint(vars["appointment_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	appointment := models.Appointment{}
	err = json.Unmarshal(body, &appointment)
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

	err = appointment.Validate("update")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedAppointment, err := appointment.UpdateAppointment(server.DB, uint32(aid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, updatedAppointment)
}

// DeleteAppointment ...
func (server *Server) DeleteAppointment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appointment := models.Appointment{}

	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	aid, err := strconv.ParseUint(vars["appointment_id"], 10, 32)
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
	_, err = appointment.DeleteAppointment(server.DB, uint32(aid))
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", aid))
	handlers.ResponseJSON(w, http.StatusNoContent, "")
}
