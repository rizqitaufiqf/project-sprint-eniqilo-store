package staff_repository

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"

	"github.com/jackc/pgx/v5"
)

type StaffRepository interface {
	Register(ctx context.Context, tx pgx.Tx, req staff_entity.Staff) (staff_entity.Staff, error)
	Login(ctx context.Context, tx pgx.Tx, req staff_entity.Staff) (staff_entity.Staff, error)
}
