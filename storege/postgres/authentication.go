package postgres

import (
	user "authentication_service/genproto/authentication_service"
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Register(userReq *user.RegisterRequest) (*user.RegisterResponse, error) {
	query := `INSERT INTO users (id, username, email, password_hash, full_name) VALUES ($1, $2, $3, $4, $5) `
	password, err := HashPassword(userReq.Password)
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()

	userReq.Password = string(password)

	_, err = r.db.Exec(query, id, userReq.Username, userReq.Email, userReq.Password, userReq.FullName)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		Id:        id,
		Username:  userReq.Username,
		Email:     userReq.Email,
		FullName:  userReq.FullName,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil

}

func (r *UserRepo) Login(userReq *user.LoginRequest) (*user.LoginResponse, error) {
	res := user.LoginResponse{}
	query :=
		`
			SELECT password_hash from users WHERE email = $1 and deleted_at is null
		`

	var password_hash string
	err := r.db.QueryRow(query, userReq.Email).Scan(&password_hash)
	if err != nil {
		if err == sql.ErrNoRows {
			res.Seccses = false
			return &res, fmt.Errorf("error getting user from database: not rows found")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(userReq.Password))
	if err != nil {
		res.Seccses = false
		return &res, fmt.Errorf("invalid password")
	}

	res.Seccses = true
	return &res, nil
}

func (r *UserRepo) GetProfile(req *user.ProfilRequest) (*user.ProfileResponse, error) {
	query := `
		SELECT id, username, email, full_name, eco_points, created_at, updated_at FROM users WHERE id = $1 and deleted_at is null 
	`
	row := r.db.QueryRow(query, req.Id)
	var res user.ProfileResponse
	err := row.Scan(&res.Id, &res.Username, &res.Email, &res.FullName, &res.EcoPoints, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil

}

func (r *UserRepo) UpdateProfile(req *user.EditProfileRequest) (*user.EditProfileResponse, error) {
	query := `
		UPDATE users SET full_name = $1, bio = $2, updated_at = now() WHERE id = $3 and deleted_at is null 
		RETURNING id, username, email, full_name, bio, updated_at
	`
	var res user.EditProfileResponse
	err := r.db.QueryRow(query, req.FullName, req.Bio, req.Id).Scan(&res.Id, &res.Username, &res.Email, &res.FullName, &res.Bio, &res.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepo) GetUsers(req *user.UsersListRequest) (*user.UsersListResponse, error) {
	query := `
		SELECT id, username, email, full_name, eco_points, created_at, updated_at FROM users WHERE deleted_at is null
	`

	// Add filtering based on request parameters
	if req.Username != "" {
		query += " AND username LIKE $1"
	}
	if req.Id != "" {
		query += " AND id LIKE $2"
	}
	if req.FullName != "" {
		query += " AND full_name LIKE $3"
	}
	if req.EcoPoints != 0 {
		query += " AND eco_points = $4"
	}
	


	rows, err := r.db.Query(query, req.Username, req.Id, req.FullName, req.EcoPoints)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*user.UsersListRequest
	for rows.Next() {
		var u user.UsersListRequest
		err := rows.Scan(&u.Id, &u.Username,  &u.FullName, &u.EcoPoints)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	// Count total users
	var total int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE deleted_at is null`).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &user.UsersListResponse{
		Users: users,
		Total: int32(total),
		Page:  (int32(total)/10 )+ 1,
		Limit: 10,
	}, nil
}


func (r *UserRepo) DeleteProfile(req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	query := `
		UPDATE users SET deleted_at = now() WHERE id = $1 and deleted_at is null 
	`
	_, err := r.db.Exec(query, req.Id)
	if err != nil {
		return &user.DeleteUserResponse{Deleted: false}, err
	}
	return &user.DeleteUserResponse{Deleted: true}, nil
}

func (r *UserRepo) ResetPassword(req *user.ResetPasswordRequest) (*user.ResetPasswordResponse, error)  {
	return nil, nil
}

func (r *UserRepo) RefreshToken(req *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error)  {
	return nil, nil
}

func (r *UserRepo) Logout(req *user.LogoutRequest) (*user.LogoutResponse, error)  {
	return nil, nil	
}

func (r *UserRepo) AddEcoPoints(req *user.AddEcoPointsRequest) (*user.AddEcoPointsResponse, error) {
	query := `
		UPDATE users SET eco_points = eco_points + $1, bio = $2, updated_at = now() WHERE id = $3 and deleted_at is null 
		RETURNING id, eco_points, bio, updated_at, 
	`
	var res user.AddEcoPointsResponse
	res.AddedPoints = req.Points
	err := r.db.QueryRow(query, req.Points, req.Reason, req.UserId).Scan(&res.UserId, &res.EcoPoints, &res.Reason, &res.Timestamp,)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepo) GetEcoPoints(req *user.EcoPointsRequest) (*user.EcoPointsResponse, error) {
	query := `
		SELECT id, eco_points,  update_at FROM users WHERE id = $1 and deleted_at is null
	`
	var res user.EcoPointsResponse
	err := r.db.QueryRow(query, req.UserId).Scan(&res.UserId, &res.EcoPoints, &res.LastUpdated) 
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *UserRepo) GetEcoPointsHistory(req *user.GetEcoPointsHistoryRequest) (*user.EcoPointsHistoryResponse, error) {
	return nil, nil
}


func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
