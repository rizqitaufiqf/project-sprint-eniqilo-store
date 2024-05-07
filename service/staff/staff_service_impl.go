package staff_service

import (
	"context"
	staff_entity "eniqilo-store/entity/staff"
	exc "eniqilo-store/exceptions"
	helpers "eniqilo-store/helpers"
	staffRep "eniqilo-store/repository/staff"
	authService "eniqilo-store/service/auth"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffServiceImpl struct {
	StaffRepository staffRep.StaffRepository
	DBPool          *pgxpool.Pool
	AuthService     authService.AuthService
	Validator       *validator.Validate
}

func NewStaffService(staffRepository staffRep.StaffRepository, dbPool *pgxpool.Pool, authService authService.AuthService, validator *validator.Validate) StaffService {
	return &StaffServiceImpl{
		StaffRepository: staffRepository,
		DBPool:          dbPool,
		AuthService:     authService,
		Validator:       validator,
	}
}

func (service *StaffServiceImpl) Register(ctx context.Context, req staff_entity.StaffRegisterRequest) (staff_entity.StaffRegisterResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return staff_entity.StaffRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	tx, err := service.DBPool.Begin(ctx)
	if err != nil {
		return staff_entity.StaffRegisterResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(ctx)

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return staff_entity.StaffRegisterResponse{}, err
	}

	staff := staff_entity.Staff{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashPassword,
	}
	staffRegistered, err := staffRep.NewStaffRepository().Register(ctx, tx, staff)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return staff_entity.StaffRegisterResponse{}, exc.ConflictException("Staff with this phone number already registered")
		}
		return staff_entity.StaffRegisterResponse{}, err
	}

	token, err := authService.NewAuthService().GenerateToken(ctx, staffRegistered.Id)
	if err != nil {
		return staff_entity.StaffRegisterResponse{}, err
	}

	return staff_entity.StaffRegisterResponse{
		Message: "Staff registered",
		Data: &staff_entity.StaffData{
			PhoneNumber: staffRegistered.PhoneNumber,
			Name:        staffRegistered.Name,
			AccessToken: token,
		},
	}, nil
}
