package staff_repository

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"

	"github.com/jackc/pgx/v5"
)

type StaffRepositoryImpl struct {
}

func NewStaffRepository() StaffRepository {
	return &StaffRepositoryImpl{}
}

func (repository *StaffRepositoryImpl) Register(ctx context.Context, tx pgx.Tx, staff staff_entity.Staff) (staff_entity.Staff, error) {
	var staffId string
	query := "INSERT INTO staffs (id, name, phone_number, password) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id"
	if err := tx.QueryRow(ctx, query, staff.Name, staff.PhoneNumber, staff.Password).Scan(&staffId); err != nil {
		return staff_entity.Staff{}, err
	}

	staff.Id = staffId
	if err := tx.Commit(ctx); err != nil {
		return staff_entity.Staff{}, err
	}
	return staff, nil
}
