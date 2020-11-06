package controllers

import "github.com/repoerna/hms_app/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login/{user}", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// Patient routes
	s.Router.HandleFunc("/patients", middlewares.SetMiddlewareJSON(s.CreatePatient)).Methods("POST")
	s.Router.HandleFunc("/patients", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPatients))).Methods("GET")
	s.Router.HandleFunc("/patients/{ssn}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPatient))).Methods("GET")
	s.Router.HandleFunc("/patients/{ssn}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePatient))).Methods("PUT")
	s.Router.HandleFunc("/patients/{ssn}", middlewares.SetMiddlewareAuthentication(s.DeletePatient)).Methods("DELETE")

	// Employee routes
	s.Router.HandleFunc("/employees", middlewares.SetMiddlewareJSON(s.CreateEmployee)).Methods("POST")
	s.Router.HandleFunc("/employees", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetEmployees))).Methods("GET")
	s.Router.HandleFunc("/employees/{employee_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetEmployee))).Methods("GET")
	s.Router.HandleFunc("/employees/{employee_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateEmployee))).Methods("PUT")
	s.Router.HandleFunc("/employees/{employee_id}", middlewares.SetMiddlewareAuthentication(s.DeleteEmployee)).Methods("DELETE")

	// Schedule routes
	s.Router.HandleFunc("/schedules", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateSchedule))).Methods("POST")
	s.Router.HandleFunc("/schedules", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetSchedules))).Methods("GET")
	s.Router.HandleFunc("/schedules/{user_id}/{schedule_code}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateSchedule))).Methods("PUT")
	s.Router.HandleFunc("/schedules/{user_id}/{schedule_code}", middlewares.SetMiddlewareAuthentication(s.DeleteSchedule)).Methods("DELETE")

	// Apointment routes
	s.Router.HandleFunc("/appointments", middlewares.SetMiddlewareJSON(s.CreateAppointment)).Methods("POST")
	s.Router.HandleFunc("/appointments", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetAppointments))).Methods("GET")
	s.Router.HandleFunc("/appointments/{user_id}/{appointment_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetAppointment))).Methods("GET")
	s.Router.HandleFunc("/appointments/{user_id}/{appointment_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateAppointment))).Methods("PUT")
	s.Router.HandleFunc("/appointments/{user_id}/{appointment_id}", middlewares.SetMiddlewareAuthentication(s.DeleteAppointment)).Methods("DELETE")

	// examinations routes
	s.Router.HandleFunc("/examinations", middlewares.SetMiddlewareJSON(s.CreateExamination)).Methods("POST")
	s.Router.HandleFunc("/examinations", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetExamination))).Methods("GET")
	s.Router.HandleFunc("/examinations/{user_id}/{examination_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetExamination))).Methods("GET")
	s.Router.HandleFunc("/examinations/{user_id}/{examination_id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateExamination))).Methods("PUT")
	s.Router.HandleFunc("/examinations/{user_id}/{examination_id}", middlewares.SetMiddlewareAuthentication(s.DeleteExamination)).Methods("DELETE")
}
