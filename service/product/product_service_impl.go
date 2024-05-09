package product_service

import (
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	product_repository "eniqilo-store/repository/product"
	"fmt"

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
			CreatedAt: productAdded.CreatedAt,
		},
	}, nil

}
