package staff_repository

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"
)

type StaffRepository interface {
	Register(ctx context.Context, req staff_entity.Staff) (staff_entity.Staff, error)
	Login(ctx context.Context, req staff_entity.Staff) (staff_entity.Staff, error)
}
