package product_service

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
	product_repository "eniqilo-store/repository/product"
	"fmt"
	"strings"

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
	productAdded, err := service.ProductRepository.Add(ctx.UserContext(), req)
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

func (service *productServiceImpl) Edit(ctx *fiber.Ctx, req product_entity.ProductEditRequest) (product_entity.ProductEditResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return product_entity.ProductEditResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err.Error()))
	}
	productId := ctx.Params("id")
	editedProduct, err := service.ProductRepository.Edit(ctx.UserContext(), req, productId)
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

	productSearched, err := service.ProductRepository.Search(ctx.UserContext(), searchQueries)
	if err != nil {
		return product_entity.ProductSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	return product_entity.ProductSearchResponse{
		Message: "List of products retrieved",
		Data:    productSearched,
	}, nil
}

func (service *productServiceImpl) Delete(ctx *fiber.Ctx) (product_entity.ProductDeleteResponse, error) {
	productId := ctx.Params("id")
	userCtx := ctx.UserContext()
	productDeleted, err := service.ProductRepository.Delete(userCtx, productId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return product_entity.ProductDeleteResponse{}, exc.NotFoundException(fmt.Sprintf("Product with id %s Not Found", productId))
		}

		return product_entity.ProductDeleteResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err.Error()))
	}

	return product_entity.ProductDeleteResponse{
		Message: "Product successfully deleted",
		Data:    productDeleted,
	}, nil
}

func (service *productServiceImpl) Checkout(ctx *fiber.Ctx, req product_entity.ProductCheckoutRequest) (product_entity.ProductCheckoutResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return product_entity.ProductCheckoutResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err.Error()))
	}

	productCheckout := product_entity.ProductCheckout{
		CustomerId:     req.CustomerId,
		ProductDetails: req.ProductDetails,
		Paid:           &req.Paid,
		Change:         req.Change,
	}

	userCtx := ctx.UserContext()
	productAdded, err := service.ProductRepository.Checkout(userCtx, productCheckout)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return product_entity.ProductCheckoutResponse{}, exc.NotFoundException("either customer or product id not found")
		}
		if strings.Contains(err.Error(), "doesn’t pass validation") {
			return product_entity.ProductCheckoutResponse{}, exc.BadRequestException(err.Error())
		}

		return product_entity.ProductCheckoutResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err.Error()))
	}

	return product_entity.ProductCheckoutResponse{
		Message: "Checkout success",
		Data: &product_entity.ProductCheckoutData{
			Id: productAdded.CheckoutId,
		},
	}, nil
}

func (service *productServiceImpl) HistorySearch(ctx *fiber.Ctx, searchQuery product_entity.ProductCheckoutHistoryRequest) (product_entity.ProductCheckoutHistoryResponse, error) {
	if err := service.Validator.Struct(searchQuery); err != nil {
		return product_entity.ProductCheckoutHistoryResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}
	if strings.ToLower(searchQuery.CreatedAt) != "asc" {
		searchQuery.CreatedAt = "desc"
	}
	historySearched, err := service.ProductRepository.HistorySearch(ctx.UserContext(), searchQuery)
	if err != nil {
		return product_entity.ProductCheckoutHistoryResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	return product_entity.ProductCheckoutHistoryResponse{
		Message: "Checkout history successfully retrieved",
		Data:    historySearched,
	}, nil
}

func (s *productServiceImpl) CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearchQuery) (product_entity.ProductCustomerSearchResponse, error) {
	if err := s.Validator.Struct(searchQuery); err != nil {
		return product_entity.ProductCustomerSearchResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	productSearched, err := s.ProductRepository.CustomerSearch(ctx, searchQuery)
	if err != nil {
		return product_entity.ProductCustomerSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	return product_entity.ProductCustomerSearchResponse{
		Message: "success",
		Data:    productSearched,
	}, nil
}
