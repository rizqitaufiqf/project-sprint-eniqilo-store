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
	"golang.org/x/crypto/bcrypt"
)

type staffServiceImpl struct {
	StaffRepository staffRep.StaffRepository
	AuthService     authService.AuthService
	Validator       *validator.Validate
}

func NewStaffService(staffRepository staffRep.StaffRepository, authService authService.AuthService, validator *validator.Validate) StaffService {
	return &staffServiceImpl{
		StaffRepository: staffRepository,
		AuthService:     authService,
		Validator:       validator,
	}
}

func (service *staffServiceImpl) Register(ctx context.Context, req staff_entity.StaffRegisterRequest) (staff_entity.StaffRegisterResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return staff_entity.StaffRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return staff_entity.StaffRegisterResponse{}, err
	}

	staff := staff_entity.Staff{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashPassword,
	}

	staffRegistered, err := service.StaffRepository.Register(ctx, staff)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return staff_entity.StaffRegisterResponse{}, exc.ConflictException("Staff with this phone number already registered")
		}
		return staff_entity.StaffRegisterResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, staffRegistered.Id)
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

func (service *staffServiceImpl) Login(ctx context.Context, req staff_entity.StaffLoginRequest) (staff_entity.StaffLoginResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return staff_entity.StaffLoginResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	staff := staff_entity.Staff{
		PhoneNumber: req.PhoneNumber,
	}

	staffLogin, err := service.StaffRepository.Login(ctx, staff)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return staff_entity.StaffLoginResponse{}, exc.NotFoundException("Staff is not found")
		}

		return staff_entity.StaffLoginResponse{}, err
	}

	if _, err = helpers.ComparePassword(staffLogin.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return staff_entity.StaffLoginResponse{}, exc.BadRequestException("Invalid password")
		}

		return staff_entity.StaffLoginResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, staffLogin.Id)
	if err != nil {
		return staff_entity.StaffLoginResponse{}, err
	}

	return staff_entity.StaffLoginResponse{
		Message: "Staff logged successfully",
		Data: &staff_entity.StaffData{
			Id:          staffLogin.Id,
			PhoneNumber: req.PhoneNumber,
			Name:        staffLogin.Name,
			AccessToken: token,
		},
	}, nil
}
