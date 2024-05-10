package customer_entity

type CustomerRegisterRequest struct {
	Name        string `json:"name" validate:"min=5,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"min=10,max=16,phoneNumber"`
}

type CustomerSearchRequest struct {
	PhoneNumber string
	Name        string
}
