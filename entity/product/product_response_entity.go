package product_entity

type ProductRegisterResponse struct {
	Message string       `json:"message"`
	Data    *ProductData `json:"data"`
}

type ProductData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ProductCustomerSearchResponse struct {
	Message string                       `json:"message"`
	Data    *[]ProductCustomerSearchData `json:"data"`
}

type ProductCustomerSearchData struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Sku       string `json:"sku"`
	Category  string `json:"category"`
	ImageUrl  string `json:"imageUrl"`
	Stock     int    `json:"stock"`
	Price     int    `json:"price"`
	Location  string `json:"location"`
	CreatedAt string `json:"createdAt"`
}
