package staff_entity

type StaffRegisterResponse struct {
	Message string     `json:"message"`
	Data    *StaffData `json:"data"`
}

type StaffLoginResponse struct {
	Message string     `json:"message"`
	Data    *StaffData `json:"data,omitempty"`
}

type StaffData struct {
	Id          string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"accessToken"`
}
