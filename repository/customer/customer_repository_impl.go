package customer_repository

import (
	"context"
	customer_entity "eniqilo-store/entity/customer"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepositoryImpl struct {
	dbPool *pgxpool.Pool
}

func NewCustomerRepository(dbPool *pgxpool.Pool) CustomerRepository {
	return &CustomerRepositoryImpl{
		dbPool: dbPool,
	}
}

func (repository *CustomerRepositoryImpl) Register(ctx context.Context, customer customer_entity.CustomerRegisterRequest) (customer_entity.CustomerData, error) {
	var customerId string
	query := "INSERT INTO customers (id, name, phone_number) VALUES (gen_random_uuid(), $1, $2) RETURNING id"
	if err := repository.dbPool.QueryRow(ctx, query, customer.Name, customer.PhoneNumber).Scan(&customerId); err != nil {
		return customer_entity.CustomerData{}, err
	}

	return customer_entity.CustomerData{UserId: customerId, Name: customer.Name, PhoneNumber: customer.PhoneNumber}, nil
}

func (repository *CustomerRepositoryImpl) Search(ctx context.Context, customer customer_entity.CustomerSearchRequest) ([]customer_entity.CustomerData, error) {
	query := `SELECT id userId, phone_number phoneNumber, name
	FROM customers
	WHERE name ~* $1 AND phone_number LIKE $2
	ORDER BY created_at DESC
	`
	rows, err := repository.dbPool.Query(ctx, query, customer.Name, fmt.Sprintf("+%s%%", customer.PhoneNumber))
	if err != nil {
		return []customer_entity.CustomerData{}, err
	}
	defer rows.Close()

	matches, err := pgx.CollectRows(rows, pgx.RowToStructByName[customer_entity.CustomerData])

	if err != nil {
		return []customer_entity.CustomerData{}, err
	}

	return matches, nil
}
