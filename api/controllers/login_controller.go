package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/repoerna/hms_app/api/utils/formaterror"
	"github.com/repoerna/hms_app/api/utils/hash"

	"github.com/repoerna/hms_app/api/auth"
	"github.com/repoerna/hms_app/api/handlers"
	"github.com/repoerna/hms_app/api/models"
	"golang.org/x/crypto/bcrypt"
)

// Login ...
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	user := vars["user"]

	switch strings.ToLower(user) {
	case "patient":
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

		err = patient.Validate("login")
		if err != nil {
			handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
			return
		}
		token, err := server.SignIn("patient", patient.Email, patient.Password)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			handlers.ResponseError(w, http.StatusUnprocessableEntity, formattedError)
			return
		}
		handlers.ResponseJSON(w, http.StatusOK, token)

	case "employee":
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

		err = employee.Validate("login")
		if err != nil {
			handlers.ResponseError(w, http.StatusUnprocessableEntity, err)
			return
		}
		token, err := server.SignIn("employee", employee.Email, employee.Password)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			handlers.ResponseError(w, http.StatusUnprocessableEntity, formattedError)
			return
		}
		handlers.ResponseJSON(w, http.StatusOK, token)
	}

}

// SignIn ...
func (server *Server) SignIn(user, email, password string) (string, error) {
	if strings.ToLower(user) == "patient" {
		var err error

		patient := models.Patient{}

		err = server.DB.Debug().Model(models.Patient{}).Where("email = ?", email).Take(&patient).Error
		if err != nil {
			return "", err
		}
		err = hash.VerifyPassword(patient.Password, password)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			return "", err
		}
		return auth.CreateToken(patient.SSN)
	} else if strings.ToLower(user) == "employee" {
		var err error

		employee := models.Employee{}

		err = server.DB.Debug().Model(models.Employee{}).Where("email = ?", email).Take(&employee).Error
		if err != nil {
			return "", err
		}
		err = hash.VerifyPassword(employee.Password, password)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			return "", err
		}
		return auth.CreateToken(employee.EmployeeID)
	}
	return "", nil
}
