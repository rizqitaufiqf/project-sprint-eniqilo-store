package product_entity

type ProductRegisterResponse struct {
	Message string       `json:"message"`
	Data    *ProductData `json:"data"`
}

type ProductData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}
