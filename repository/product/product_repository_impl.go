package product_repository

import (
	"context"
	product_entity "eniqilo-store/entity/product"
	exc "eniqilo-store/exceptions"
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

func (repository *productRepositoryImpl) Edit(ctx context.Context, product product_entity.Product, productId string) (*product_entity.Product, error) {
	updateQ := `UPDATE products SET name = $1, sku = $2, category = $3,
	image_url = $4, notes = $5, price = $6, stock = $7, location = $8,
	is_available = $9
	WHERE id = $10
	`
	res, err := repository.dbPool.Exec(ctx, updateQ, product.Name, product.Sku, product.Category, product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, productId)
	if err != nil {
		return &product_entity.Product{}, err
	}
	if res.RowsAffected() == 0 {
		return &product_entity.Product{}, exc.NotFoundException("Product id does not exist")
	}

	product.Id = productId

	return &product, nil
}

func (repository *productRepositoryImpl) Search(ctx context.Context, searchQuery product_entity.ProductSearch) (*[]product_entity.Product, error) {
	query := `SELECT id, name, sku, category, image_url, stock, notes, price, location, is_available, created_at FROM products WHERE is_deleted = false`
	var whereClause []string
	var searchParams []interface{}

	if searchQuery.Id != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%s", strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, searchQuery.Id)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%s", strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.IsAvailable != "" {
		isAvail, err := strconv.ParseBool(searchQuery.IsAvailable)
		if err != nil {
			return &[]product_entity.Product{}, err
		}
		whereClause = append(whereClause, fmt.Sprintf("is_available = $%s", strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, isAvail)
	}
	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%s", strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, searchQuery.Category)
	}
	if searchQuery.Sku != "" {
		whereClause = append(whereClause, fmt.Sprintf("sku = $%s", strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, searchQuery.Sku)
	}
	if searchQuery.InStock != "" {
		inStock, err := strconv.ParseBool(searchQuery.InStock)
		if err != nil {
			return &[]product_entity.Product{}, err
		}
		var op string
		if inStock {
			op = ">"
		} else {
			op = "="
		}
		whereClause = append(whereClause, fmt.Sprintf("stock %s $%s", op, strconv.Itoa(len(searchParams)+1)))
		searchParams = append(searchParams, 0)
	}
	if len(whereClause) > 0 {
		query += " AND " + strings.Join(whereClause, " AND ")
	}

	// construct order by
	var orderByClause []string
	var orderByDefault = ` ORDER BY created_at DESC`
	if searchQuery.Price != "" {
		orderByClause = append(orderByClause, fmt.Sprintf("price %s", searchQuery.Price))
	}
	if searchQuery.CreatedAt != "" {
		orderByClause = append(orderByClause, fmt.Sprintf("created_at %s", searchQuery.CreatedAt))
	}

	if len(orderByClause) > 0 {
		query += " ORDER BY " + strings.Join(orderByClause, ", ")
	} else {
		query += orderByDefault
	}

	if searchQuery.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(len(searchParams)+1), strconv.Itoa(len(searchParams)+2))
		searchParams = append(searchParams, searchQuery.Limit, searchQuery.Offset)
	}

	rows, err := repository.dbPool.Query(ctx, query, searchParams...)
	if err != nil {
		return &[]product_entity.Product{}, err
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product_entity.Product])
	if err != nil {
		return &[]product_entity.Product{}, err
	}

	return &products, nil

}
