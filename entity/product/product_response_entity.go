package product_entity

type ProductRegisterResponse struct {
	Message string       `json:"message"`
	Data    *ProductData `json:"data"`
}

type ProductEditResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ProductData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ProductSearchResponse struct {
	Message string               `json:"message"`
	Data    *[]ProductSearchData `json:"data"`
}

type ProductSearchData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Sku         string `json:"sku"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	Stock       int    `json:"stock"`
	Notes       string `json:"notes"`
	Price       int    `json:"price"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedAt   string `json:"createdAt"`
}
