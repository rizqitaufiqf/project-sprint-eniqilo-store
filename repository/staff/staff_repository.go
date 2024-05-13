package staff_repository

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"
)

type StaffRepository interface {
	Register(ctx context.Context, req staff_entity.StaffRegisterRequest) (staff_entity.StaffData, error)
	Login(ctx context.Context, req staff_entity.StaffLoginRequest) (staff_entity.Staff, error)
}
