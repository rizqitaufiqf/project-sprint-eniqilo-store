package product_entity

type ProductRegisterRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Sku         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,productCategory"`
	ImageUrl    string `json:"imageUrl" validate:"required,url"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int    `json:"price" validate:"required,min=1"`
	Stock       *int   `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"required"`
}

type ProductCustomerSearchQuery struct {
	Name     string `query:"name" validate:"omitempty"`
	Category string `query:"category" validate:"omitempty"`
	Sku      string `query:"sku" validate:"omitempty"`
	Price    string `query:"price" validate:"omitempty"`
	InStock  string `query:"inStock" validate:"omitempty"`
	Limit    string `query:"limit" validate:"omitempty,number,min=0"`
	Offset   string `query:"offset" validate:"omitempty,number,min=0"`
}
