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

type ProductCheckoutDetailsRequest struct {
	ProductId string `json:"productId" validate:"required,min=1"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type ProductCheckoutRequest struct {
	CustomerId     string                           `json:"customerId" validate:"required,min=1"`
	ProductDetails *[]ProductCheckoutDetailsRequest `json:"productDetails" validate:"required,min=1"`
	Paid           int                              `json:"paid" validate:"required,min=1"`
	Change         int                              `json:"change" validate:"required,min=1"`
}
