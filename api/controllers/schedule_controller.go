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

// CreateSchedule ...
func (server *Server) CreateSchedule(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
	}
	schedule := models.Schedule{}
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = schedule.Validate("")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	scheduleCreated, err := schedule.SaveSchedule(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, scheduleCreated.ScheduleCode))
	handlers.ResponseJSON(w, http.StatusCreated, scheduleCreated)
}

// GetSchedules ...
func (server *Server) GetSchedules(w http.ResponseWriter, r *http.Request) {

	schedule := models.Schedule{}

	schedules, err := schedule.FindAllSchedules(server.DB)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, schedules)
}

// UpdateSchedule ...
func (server *Server) UpdateSchedule(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)

	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	sc := vars["schedule_code"]
	if sc == "" {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	schedule := models.Schedule{}
	err = json.Unmarshal(body, &schedule)
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

	err = schedule.Validate("update")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedSchedule, err := schedule.UpdateSchedule(server.DB, sc)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, updatedSchedule)
}

// DeleteSchedule ...
func (server *Server) DeleteSchedule(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	schedule := models.Schedule{}

	uid, err := strconv.ParseUint(vars["user_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	sc := vars["schedule_code"]
	if sc == "" {
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
	_, err = schedule.DeleteSchedule(server.DB, sc)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", sc))
	handlers.ResponseJSON(w, http.StatusNoContent, "")
}
