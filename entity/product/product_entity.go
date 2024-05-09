package product_entity

type Product struct {
	Id, Name, Sku, Category, ImageUrl, Notes, Location string
	Price                                              int
	Stock                                              int
	IsAvailable, IsDeleted                             bool
	CreatedAt                                          string
}

type ProductCheckout struct {
	CheckoutId, CustomerId string
	Paid, Change           int
	ProductDetails         *[]ProductCheckoutDetailsRequest
	CreatedAt              string
}
