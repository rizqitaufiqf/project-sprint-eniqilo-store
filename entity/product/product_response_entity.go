package product_entity

type ProductRegisterResponse struct {
	Message string       `json:"message"`
	Data    *ProductData `json:"data"`
}

type ProductData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ProductDeleteResponse struct {
	Message string             `json:"message"`
	Data    *ProductDeleteData `json:"data"`
}

type ProductDeleteData struct {
	Id string `json:"id"`
}
