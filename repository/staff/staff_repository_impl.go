package staff_repository

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"

	"github.com/jackc/pgx/v5/pgxpool"
)

type staffRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewStaffRepository(dbPool *pgxpool.Pool) StaffRepository {
	return &staffRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *staffRepositoryImpl) Register(ctx context.Context, staff staff_entity.Staff) (staff_entity.Staff, error) {
	var staffId string
	query := "INSERT INTO staffs (id, name, phone_number, password) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, staff.Name, staff.PhoneNumber, staff.Password).Scan(&staffId); err != nil {
		return staff_entity.Staff{}, err
	}

	staff.Id = staffId
	return staff, nil
}

func (repository *staffRepositoryImpl) Login(ctx context.Context, staff staff_entity.Staff) (staff_entity.Staff, error) {
	query := "SELECT id, name, password FROM staffs WHERE phone_number = $1 LIMIT 1"
	row := repository.DBpool.QueryRow(ctx, query, staff.PhoneNumber)

	var loggedInStaff staff_entity.Staff
	err := row.Scan(&loggedInStaff.Id, &loggedInStaff.Name, &loggedInStaff.Password)
	if err != nil {
		return staff_entity.Staff{}, err
	}

	return loggedInStaff, nil
}
