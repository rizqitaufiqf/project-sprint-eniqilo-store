package product_entity

type Product struct {
	Id, Name, Sku, Category, ImageUrl, Notes, Location string
	Price                                              int
	Stock                                              int
	IsAvailable                                        bool
	CreatedAt                                          string
}

type ProductSearch struct {
	Id          string
	Name        string
	IsAvailable string
	Category    string
	Sku         string
	Price       int
	InStock     string
	CreatedAt   string
	Limit       int
	Offset      int
}
