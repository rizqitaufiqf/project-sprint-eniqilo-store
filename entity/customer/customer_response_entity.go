package customer_entity

type CustomerRegisterResponse struct {
	Message string        `json:"message"`
	Data    *CustomerData `json:"data"`
}

type CustomerSearchResponse struct {
	Message string          `json:"message"`
	Data    *[]CustomerData `json:"data,omitempty"`
}

type CustomerData struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
