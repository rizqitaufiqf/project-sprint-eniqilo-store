package staff_entity

type StaffResponse struct {
	Message string     `json:"message"`
	Data    *StaffData `json:"data"`
}

type StaffData struct {
	Id          string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"accessToken"`
}
