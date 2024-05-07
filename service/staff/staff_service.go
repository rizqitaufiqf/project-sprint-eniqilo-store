package staff_service

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"
)

type StaffService interface {
	Register(ctx context.Context, req staff_entity.StaffRegisterRequest) (staff_entity.StaffRegisterResponse, error)
	Login(ctx context.Context, req staff_entity.StaffLoginRequest) (staff_entity.StaffLoginResponse, error)
}
