package product_entity

type Product struct {
	Id, Name, Sku, Category, ImageUrl, Notes, Location string
	Price                                              int
	Stock                                              int
	IsAvailable                                        bool
	CreatedAt                                          string
}
