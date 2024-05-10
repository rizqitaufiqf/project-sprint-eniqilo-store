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

type ProductCheckoutData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ProductCheckoutResponse struct {
	Message string               `json:"message"`
	Data    *ProductCheckoutData `json:"data"`
}

type CheckoutDetailsData struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type ProductCheckoutDataResponse struct {
	TransactionId  string                 `json:"transactionId"`
	CustomerId     string                 `json:"customerId"`
	ProductDetails *[]CheckoutDetailsData `json:"productDetails"`
	Paid           int                    `json:"paid"`
	Change         int                    `json:"change"`
	CreatedAt      string                 `json:"createdAt"`
}

type ProductCheckoutHistoryResponse struct {
	Message string                         `json:"message"`
	Data    *[]ProductCheckoutDataResponse `json:"data"`
}
