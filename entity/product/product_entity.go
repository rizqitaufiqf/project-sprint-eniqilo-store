package product_entity

import "time"

type Product struct {
	Id, Name, Sku, Category, ImageUrl, Notes, Location string
	Price                                              int
	Stock                                              int
	IsAvailable, IsDeleted                             bool
	CreatedAt                                          time.Time
}

type ProductCustomerSearch struct {
	Name, Sku, Category, InStock, Price string
	Limit, Offset                       int
}

type ProductCheckout struct {
	CheckoutId, CustomerId string
	Paid, Change           *int
	ProductDetails         *[]ProductCheckoutDetailsRequest
	CreatedAt              string
}
