package models

import ()

type LoginRequestDevice struct {
	OsType		string	`json:"os_type"`
	DevicePush	bool	`json:"device_push"`
	Platform	string	`json:"platform"`
	UserAgent	string	`json:"user_agent"`
	Vendor		string	`json:"vendor"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Device	 LoginRequestDevice	`json:"device"`
}

type RawRegistrationRequest struct {
	Data map[string]interface{} `json:"data"`
}

type RegistrationRequest struct {
	Data RegistrationData `json:"data"`
}

type RegistrationData struct {
	FirstName         string    `json:"FirstName"`
	LastName          string    `json:"LastName"`
	Email             string    `json:"EmailAddress"`
	Password          string    `json:"Password"`
	DateOfBirth       string    `json:"DateOfBirth"`
	InsuranceID       string    `json:"InsuranceID,omitempty"`
	EmployerName      string    `json:"EmployerCompanyName,omitempty"`
	DisplayName       string    `json:"NewtopiaProfileName,omitempty"`
	PhotoType         string    `json:"ProfilePhoto,omitempty"`
	PhotoUpload       string    `json:"UploadPhoto,omitempty"`  // ProfilePhoto = 0
	PhotoSelect       string    `json:"SelectAvatar,omitempty"` // ProfilePhoto = 1
	CommunityGreeting string    `json:"HelloNewtopiaCommunity,omitempty"`
	AddressLine1      string    `json:"AddressLine1,omitempty"`
	AddressLine2      string    `json:"AddressLine2,omitempty"`
	City              SelectBox `json:"AddressCity,omitempty"`
	State             SelectBox `json:"AddressState,omitempty"`
	Phone             string    `json:"TODO,omitempty"`
	ConsentGenetic    Checkbox  `json:"ExpressConsentGenetic"`
	ConsentCSA        Checkbox  `json:"ExpressConsentCSA"`
}

type SelectBox struct {
	Key   string `json:"payload"`
	Value string `json:"text"`
}

type Checkbox struct {
	Value bool `json:"0"`
}
