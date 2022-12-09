package authorization

import (
	"blockpost/config"
	"blockpost/genprotos/authorization"
	"blockpost/storage"
	"blockpost/util"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// authorizationService is a struct that implements the server interface
type authorizationService struct {
	stg storage.StorageI
	cfg config.Config
	authorization.UnimplementedAuthServiceServer
}

//NewAuthService...
func NewAuthService(cfg config.Config, stg storage.StorageI) *authorizationService {
	return &authorizationService{
		cfg: cfg,
		stg: stg,
	}
}

// CreateUser ...
func (a *authorizationService) CreateUser(c context.Context, req *authorization.CreateUserRequest) (*authorization.User, error) {
	id := uuid.New()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword

	err = a.stg.CreateUser(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateUser: %s", err.Error())
	}
	user, err := a.stg.GetUserByID(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "a.stg.GetUserByID: %s", err.Error())
	}
	return user, nil
}

// GetUserByID ...
func (a *authorizationService) GetUserByID(c context.Context, req *authorization.GetUserByIDRequest) (*authorization.User, error) {

	res, err := a.stg.GetUserByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetUserByID: %s", err.Error())
	}
	return res, nil
}

// GetUserByUsername ...
func (a *authorizationService) GetUserByUsername(c context.Context, req *authorization.User) (*authorization.User, error) {

	res, err := a.stg.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetUserByUsername: %s", err.Error())
	}
	return res, nil
}

// GetUserList ...
func (a *authorizationService) GetUserList(c context.Context, req *authorization.GetUserListRequest) (*authorization.GetUserListResponse, error) {
	res, err := a.stg.GetUserList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetUserList: %s", err.Error())
	}
	return res, nil
}

// UpdateUser ...
func (a *authorizationService) UpdateUser(c context.Context, req *authorization.UpdateUserRequest) (*authorization.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword
	err = a.stg.UpdateUser(req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.UpdateUser: %s", err.Error())
	}
	res, e := a.stg.GetUserByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetUserByID: %s", e.Error())
	}
	return res, nil
}

// DeleteUser ...
func (a *authorizationService) DeleteUser(c context.Context, req *authorization.DeleteUserRequest) (*authorization.User, error) {
	res, e := a.stg.GetUserByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.GetUserByID: %s", e.Error())
	}
	err := a.stg.DeleteUser(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.stg.DeleteUser: %s", err.Error())
	}

	return res, nil
}
