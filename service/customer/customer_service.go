package customer_service

import (
	"context"
	customer_entity "eniqilo-store/entity/customer"
)

type CustomerService interface {
	Register(ctx context.Context, req customer_entity.CustomerRegisterRequest) (customer_entity.CustomerRegisterResponse, error)
	Search(ctx context.Context, req customer_entity.CustomerSearchRequest) (customer_entity.CustomerSearchResponse, error)
}
