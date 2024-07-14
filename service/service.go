package service

import (
	pb "authentication_service/genproto/authentication_service"
	"authentication_service/storege/postgres"
	"context"
)


type AuthService struct {
	pb.UnimplementedEcoServiceServer
	repo *postgres.UserRepo
}

func NewAuthService(db *postgres.UserRepo) *AuthService {
	return &AuthService{repo: db}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.repo.Register(req)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.repo.Login(req)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *AuthService) GetProfile(ctx context.Context, UserId *pb.ProfilRequest) (*pb.ProfileResponse, error) {
	user, err := s.repo.GetProfile(UserId)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *AuthService) EditProfile(ctx context.Context, req *pb.EditProfileRequest) (*pb.EditProfileResponse, error)  {
	user, err := s.repo.UpdateProfile(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ListUsers(ctx context.Context, req *pb.UsersListRequest) (*pb.UsersListResponse, error) {
	users, err := s.repo.GetUsers(req)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user, err := s.repo.DeleteProfile(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error)  {
	return nil, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)  {
	return nil, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)  {
	return nil, nil	
}

func (s *AuthService) AddEcoPoints(ctx context.Context, req *pb.AddEcoPointsRequest) (*pb.AddEcoPointsResponse, error) {
	user, err := s.repo.AddEcoPoints(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetEcoPoints(ctx context.Context, req *pb.EcoPointsRequest) (*pb.EcoPointsResponse, error) {
	user, err := s.repo.GetEcoPoints(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetEcoPointsHistory(ctx context.Context, req *pb.GetEcoPointsHistoryRequest) (*pb.EcoPointsHistoryResponse, error) {
	return nil, nil
}


