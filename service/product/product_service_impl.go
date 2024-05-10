package product_service

import (
	"context"
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

func (s *productServiceImpl) CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearchQuery) (product_entity.ProductCustomerSearchResponse, error) {
	if err := s.Validator.Struct(searchQuery); err != nil {
		return product_entity.ProductCustomerSearchResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	if strings.ToLower(searchQuery.InStock) != "true" && strings.ToLower(searchQuery.InStock) != "false" {
		searchQuery.InStock = ""
	}

	if strings.ToLower(searchQuery.Price) != "asc" && strings.ToLower(searchQuery.Price) != "desc" {
		searchQuery.Price = ""
	}

	exist := false
	for _, category := range helpers.ProductCategory {
		if strings.EqualFold(category, searchQuery.Category) {
			exist = true
			break
		}
	}

	if !exist {
		searchQuery.Category = ""
	}

	productQuery := product_entity.ProductCustomerSearch{
		Name:     searchQuery.Name,
		Sku:      searchQuery.Sku,
		Category: searchQuery.Category,
		Price:    searchQuery.Price,
		InStock:  searchQuery.InStock,
		Limit:    5,
		Offset:   0,
	}
	if searchQuery.Limit != "" {
		productQuery.Limit, _ = strconv.Atoi(searchQuery.Limit)
	}
	if searchQuery.Offset != "" {
		productQuery.Offset, _ = strconv.Atoi(searchQuery.Offset)
	}

	productSearched, err := s.ProductRepository.CustomerSearch(ctx, productQuery)
	if err != nil {
		return product_entity.ProductCustomerSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	return product_entity.ProductCustomerSearchResponse{
		Message: "success",
		Data:    productSearched,
	}, nil
}
