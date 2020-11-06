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

// CreateEmployee ...
func (server *Server) CreateEmployee(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
	}
	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = employee.Validate("")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	employeeCreated, err := employee.SaveEmployee(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, employeeCreated.EmployeeID))
	handlers.ResponseJSON(w, http.StatusCreated, employeeCreated)
}

// GetEmployees ...
func (server *Server) GetEmployees(w http.ResponseWriter, r *http.Request) {

	employee := models.Employee{}

	employees, err := employee.FindAllEmployee(server.DB)
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, employees)
}

// GetEmployee ...
func (server *Server) GetEmployee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	employeeID, err := strconv.ParseUint(vars["employee_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	employee := models.Employee{}
	employeeGotten, err := employee.FindEmployeeByID(server.DB, int(employeeID))
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, employeeGotten)
}

// UpdateEmployee ...
func (server *Server) UpdateEmployee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eployeeID, err := strconv.ParseUint(vars["employee_id"], 10, 32)
	if err != nil {
		handlers.ResponseError(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(eployeeID) {
		handlers.ResponseError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = employee.Validate("update")
	if err != nil {
		handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := employee.UpdateEmployee(server.DB, uint32(eployeeID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handlers.ResponseError(w, http.StatusInternalServerError, formattedError)
		return
	}
	handlers.ResponseJSON(w, http.StatusOK, updatedUser)
}

// DeleteEmployee ...
func (server *Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	employee := models.Employee{}

	uid, err := strconv.ParseUint(vars["employee_id"], 10, 32)
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
	_, err = employee.DeleteEmployee(server.DB, uint32(uid))
	if err != nil {
		handlers.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	handlers.ResponseJSON(w, http.StatusNoContent, "")
}
