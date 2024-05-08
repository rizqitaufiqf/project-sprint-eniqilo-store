package customer_repository

import (
	"context"
	customer_entity "eniqilo-store/entity/customer"
)

type CustomerRepository interface {
	Register(ctx context.Context, req customer_entity.CustomerRegisterRequest) (customer_entity.CustomerData, error)
	Search(ctx context.Context, req customer_entity.CustomerSearchRequest) ([]customer_entity.CustomerData, error)
}
