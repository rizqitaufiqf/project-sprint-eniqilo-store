package staff_entity

type StaffRegisterRequest struct {
	Name        string `json:"name" validate:"required,min=5,max=50"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16,phoneNumber"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type StaffLoginRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10,max=16,phoneNumber"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}
