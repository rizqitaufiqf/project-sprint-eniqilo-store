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

func (service *staffServiceImpl) Register(ctx context.Context, req staff_entity.StaffRegisterRequest) (staff_entity.StaffResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return staff_entity.StaffResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return staff_entity.StaffResponse{}, err
	}

	req.Password = hashPassword

	staffRegistered, err := service.StaffRepository.Register(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return staff_entity.StaffResponse{}, exc.ConflictException("Staff with this phone number already registered")
		}
		return staff_entity.StaffResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, staffRegistered.Id)
	if err != nil {
		return staff_entity.StaffResponse{}, err
	}

	staffRegistered.AccessToken = token
	return staff_entity.StaffResponse{
		Message: "Staff registered",
		Data:    &staffRegistered,
	}, nil
}

func (service *staffServiceImpl) Login(ctx context.Context, req staff_entity.StaffLoginRequest) (staff_entity.StaffResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return staff_entity.StaffResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	staffLogin, err := service.StaffRepository.Login(ctx, req)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return staff_entity.StaffResponse{}, exc.NotFoundException("Staff is not found")
		}

		return staff_entity.StaffResponse{}, err
	}

	if _, err = helpers.ComparePassword(staffLogin.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return staff_entity.StaffResponse{}, exc.BadRequestException("Invalid password")
		}

		return staff_entity.StaffResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(ctx, staffLogin.Id)
	if err != nil {
		return staff_entity.StaffResponse{}, err
	}

	return staff_entity.StaffResponse{
		Message: "Staff logged successfully",
		Data: &staff_entity.StaffData{
			Id:          staffLogin.Id,
			PhoneNumber: req.PhoneNumber,
			Name:        staffLogin.Name,
			AccessToken: token,
		},
	}, nil
}
