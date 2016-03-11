package models

import ()

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type RegistrationResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type RedirectResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"redirect_url"`
	Message string `json:"message"`
}
