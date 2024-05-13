package customer_service

import (
	"context"
	customer_entity "eniqilo-store/entity/customer"
	exc "eniqilo-store/exceptions"
	customerRep "eniqilo-store/repository/customer"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type CustomerServiceImpl struct {
	CustomerRepository customerRep.CustomerRepository
	Validator          *validator.Validate
}

func NewCustomerService(customerRepository customerRep.CustomerRepository, validator *validator.Validate) CustomerService {
	return &CustomerServiceImpl{
		CustomerRepository: customerRepository,
		Validator:          validator,
	}
}

func (service *CustomerServiceImpl) Register(ctx context.Context, req customer_entity.CustomerRegisterRequest) (customer_entity.CustomerRegisterResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return customer_entity.CustomerRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	customerRegistered, err := service.CustomerRepository.Register(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return customer_entity.CustomerRegisterResponse{}, exc.ConflictException("Customer with this phone number already registered")
		}
		return customer_entity.CustomerRegisterResponse{}, err
	}

	customerRegistered.Name = req.Name
	customerRegistered.PhoneNumber = req.PhoneNumber
	return customer_entity.CustomerRegisterResponse{
		Message: "Customer registered",
		Data:    &customerRegistered,
	}, nil
}

func (service *CustomerServiceImpl) Search(ctx context.Context, req customer_entity.CustomerSearchRequest) (customer_entity.CustomerSearchResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return customer_entity.CustomerSearchResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	customerSearch, err := service.CustomerRepository.Search(ctx, req)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return customer_entity.CustomerSearchResponse{}, exc.NotFoundException("Customer is not found")
		}

		return customer_entity.CustomerSearchResponse{}, err
	}

	return customer_entity.CustomerSearchResponse{
		Message: "Successfully retrieved customers",
		Data:    &customerSearch,
	}, nil

}
