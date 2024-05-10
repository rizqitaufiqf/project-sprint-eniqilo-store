package product_entity

type ProductRegisterRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Sku         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,productCategory"`
	ImageUrl    string `json:"imageUrl" validate:"required,validateUrl"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int    `json:"price" validate:"required,min=1"`
	Stock       *int   `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"required"`
}

type ProductEditRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Sku         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,productCategory"`
	ImageUrl    string `json:"imageUrl" validate:"required,validateUrl"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int    `json:"price" validate:"required,min=1"`
	Stock       *int   `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"required"`
}

type ProductSearchQuery struct {
	Id          string `query:"id" validate:"omitempty"`
	Name        string `query:"name" validate:"omitempty"`
	IsAvailable string `query:"isAvailable" validate:"omitempty"`
	Category    string `query:"category" validate:"omitempty"`
	Sku         string `query:"sku" validate:"omitempty"`
	Price       string `query:"price" validate:"omitempty"`
	InStock     string `query:"inStock" validate:"omitempty"`
	CreatedAt   string `query:"createdAt" validate:"omitempty"`
	Limit       string `query:"limit" validate:"omitempty,number,min=0"`
	Offset      string `query:"offset" validate:"omitempty,number,min=0"`
}

type ProductCheckoutDetailsRequest struct {
	ProductId string `json:"productId" validate:"required,min=1"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type ProductCheckoutRequest struct {
	CustomerId     string                           `json:"customerId" validate:"required,min=1"`
	ProductDetails *[]ProductCheckoutDetailsRequest `json:"productDetails" validate:"required,min=1"`
	Paid           int                              `json:"paid" validate:"required,min=1"`
	Change         *int                             `json:"change" validate:"required,min=0"`
}

type ProductCheckoutHistoryRequest struct {
	CustomerId string `query:"customerId"`
	Limit      int    `query:"limit" validate:"omitempty,number,min=0"`
	Offset     int    `query:"offset" validate:"omitempty,number,min=0"`
	CreatedAt  string `query:"createdAt "`
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
