package product_service

import (
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	"eniqilo-store/helpers"
	product_repository "eniqilo-store/repository/product"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type productServiceImpl struct {
	ProductRepository product_repository.ProductRepository
	Validator         *validator.Validate
}

func NewProductService(productRepository product_repository.ProductRepository, validator *validator.Validate) ProductService {
	return &productServiceImpl{
		ProductRepository: productRepository,
		Validator:         validator,
	}
}

func (service *productServiceImpl) Add(ctx *fiber.Ctx, req product_entity.ProductRegisterRequest) (product_entity.ProductRegisterResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return product_entity.ProductRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err.Error()))
	}

	product := product_entity.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		Category:    req.Category,
		ImageUrl:    req.ImageUrl,
		Notes:       req.Notes,
		Price:       req.Price,
		Stock:       *req.Stock,
		Location:    req.Location,
		IsAvailable: *req.IsAvailable,
	}

	userCtx := ctx.UserContext()
	productAdded, err := service.ProductRepository.Add(userCtx, product)
	if err != nil {
		return product_entity.ProductRegisterResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err.Error()))
	}

	return product_entity.ProductRegisterResponse{
		Message: "Product successfully added",
		Data: &product_entity.ProductData{
			Id:        productAdded.Id,
			CreatedAt: productAdded.CreatedAt.Format(time.RFC3339),
		},
	}, nil

}

func (service *productServiceImpl) Edit(ctx *fiber.Ctx, req product_entity.ProductEditRequest) (product_entity.ProductEditResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return product_entity.ProductEditResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err.Error()))
	}
	productId := ctx.Params("id")
	userCtx := ctx.UserContext()

	product := product_entity.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		Category:    req.Category,
		ImageUrl:    req.ImageUrl,
		Notes:       req.Notes,
		Price:       req.Price,
		Stock:       *req.Stock,
		Location:    req.Location,
		IsAvailable: *req.IsAvailable,
	}

	editedProduct, err := service.ProductRepository.Edit(userCtx, product, productId)
	if err != nil {
		return product_entity.ProductEditResponse{}, err
	}

	return product_entity.ProductEditResponse{
		Message: "Sucess edit product",
		Id:      editedProduct.Id,
	}, nil
}

func (service *productServiceImpl) Search(ctx *fiber.Ctx, searchQueries product_entity.ProductSearchQuery) (product_entity.ProductSearchResponse, error) {
	if err := service.Validator.Struct(searchQueries); err != nil {
		return product_entity.ProductSearchResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err.Error()))
	}

	userCtx := ctx.UserContext()
	isAvail := strings.ToLower(searchQueries.IsAvailable)
	inStock := strings.ToLower(searchQueries.InStock)
	if isAvail != "true" || isAvail != "false" {
		searchQueries.IsAvailable = ""
	}
	if inStock != "true" || inStock != "false" {
		searchQueries.InStock = ""
	}

	createdAt := strings.ToLower(searchQueries.CreatedAt)
	price := strings.ToLower(searchQueries.Price)
	if price != "asc" || price != "desc" {
		searchQueries.Price = ""
	}
	if createdAt != "asc" || createdAt != "desc" {
		searchQueries.CreatedAt = ""
	}

	var validCategory bool
	for _, categ := range helpers.ProductCategory {
		if categ == searchQueries.Category {
			validCategory = true
		}
	}
	if !validCategory {
		searchQueries.Category = ""
	}

	product := product_entity.ProductSearch{
		Id:          searchQueries.Id,
		Name:        searchQueries.Name,
		IsAvailable: searchQueries.IsAvailable,
		Category:    searchQueries.Category,
		Sku:         searchQueries.Sku,
		Price:       searchQueries.Price,
		InStock:     searchQueries.InStock,
		CreatedAt:   searchQueries.CreatedAt,
		Limit:       5,
		Offset:      0,
	}

	if searchQueries.Limit != "" {
		product.Limit, _ = strconv.Atoi(searchQueries.Limit)
	}
	if searchQueries.Offset != "" {
		product.Offset, _ = strconv.Atoi(searchQueries.Offset)
	}

	productSearched, err := service.ProductRepository.Search(userCtx, product)
	if err != nil {
		return product_entity.ProductSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	data := []product_entity.ProductSearchData{}
	for _, product := range *productSearched {
		productRecord := product_entity.ProductSearchData{
			Id:          product.Id,
			Name:        product.Name,
			Sku:         product.Sku,
			Category:    product.Category,
			ImageUrl:    product.ImageUrl,
			Stock:       product.Stock,
			Notes:       product.Notes,
			Price:       product.Price,
			Location:    product.Location,
			IsAvailable: product.IsAvailable,
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		}
		data = append(data, productRecord)
	}

	return product_entity.ProductSearchResponse{
		Message: "List of products retrieved",
		Data:    &data,
	}, nil

}
