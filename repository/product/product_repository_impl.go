package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepositoryImpl struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) ProductRepository {
	return &productRepositoryImpl{
		dbPool: dbPool,
	}
}

func (repository *productRepositoryImpl) Add(ctx context.Context, product product_entity.Product) (*product_entity.Product, error) {
	var productId string
	var createdAt time.Time
	query := `INSERT INTO products
	(
		id, name, sku, category, image_url, notes, price, stock, location, is_available
	)
	VALUES
	(
		gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	RETURNING id, created_at
	`
	if err := repository.dbPool.QueryRow(ctx, query, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable).Scan(&productId, &createdAt); err != nil {
		return &product_entity.Product{}, err
	}

	product.Id = productId
	product.CreatedAt = createdAt

	return &product, nil
}

func (repository *productRepositoryImpl) CustomerSearch(ctx context.Context, searchQuery product_entity.ProductCustomerSearch) (*[]product_entity.ProductCustomerSearchData, error) {
	query := `SELECT id, name, sku, category, image_url, stock, price, location, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') created_at FROM products WHERE is_deleted = FALSE AND is_available = TRUE`
	var whereClause []string
	params := []interface{}{}

	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Name)
	}

	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Category)
	}

	if searchQuery.Sku != "" {
		whereClause = append(whereClause, fmt.Sprintf("sku = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Sku)
	}

	if searchQuery.InStock != "" {
		inStock, _ := strconv.ParseBool(searchQuery.InStock)

		var operator string
		if inStock {
			operator = ">"
		} else {
			operator = "="
		}
		whereClause = append(whereClause, fmt.Sprintf("stock %s 0", operator))
	}

	if len(whereClause) > 0 {
		query += " AND " + strings.Join(whereClause, " AND ")
	}

	if searchQuery.Price != "" {
		query += fmt.Sprintf(" ORDER BY price %s", searchQuery.Price)
	} else {
		query += " ORDER BY created_at DESC"
	}

	if searchQuery.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(len(params)+1), strconv.Itoa(len(params)+2))
		params = append(params, searchQuery.Limit)
		params = append(params, searchQuery.Offset)
	}

	rows, err := repository.dbPool.Query(ctx, query, params...)
	if err != nil {
		return &[]product_entity.ProductCustomerSearchData{}, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.ProductCustomerSearchData])
	if err != nil {
		return &[]product_entity.ProductCustomerSearchData{}, err
	}

	return &products, nil
}
