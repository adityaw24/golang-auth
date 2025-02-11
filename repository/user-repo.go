package repository

import (
	"context"

	"github.com/adityaw24/golang-auth/model"
	"github.com/adityaw24/golang-auth/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TODO: search how to return a populated struct from function, if struct is defined

type UserRepo interface {
	CreateUser(ctx context.Context, tx *sqlx.Tx, u model.User) (string, error)
	FindByEmail(ctx context.Context, email string) (model.UserResponse, error)
	GetUserList(ctx context.Context) ([]model.User, error)
	GetUserDetail(ctx context.Context, id uuid.UUID) (model.UserDetailModel, error)
}

type userConnection struct {
	connection *sqlx.DB
}

func NewUserRepo(dbConn *sqlx.DB) UserRepo {
	return &userConnection{
		connection: dbConn,
	}
}

func (db *userConnection) GetUserList(ctx context.Context) ([]model.User, error) {
	users := make([]model.User, 0)

	query := `SELECT * FROM users`
	rows, err := db.connection.QueryxContext(ctx, query)

	if err != nil {
		utils.LogError("Repo", "func GetUserList", err)
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.User_id,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Phone,
		)

		if err != nil {
			utils.LogError("Repo", "GetUserList scan data", err)
			return users, err
		}
		users = append(users, user)
	}

	utils.CloseDB(rows)

	return users, err
}

func (db *userConnection) GetUserDetail(ctx context.Context, id uuid.UUID) (model.UserDetailModel, error) {

	var (
		userDetail model.UserDetailModel
	)

	//SQL Query
	query := `
		SELECT 
			u.user_id, u.username, u.name, u.phone, u.email 
		FROM users AS u
		WHERE user_id=$1`

	err := db.connection.QueryRowxContext(
		ctx,
		query,
		id,
	).Scan(
		&userDetail.User_id,
		&userDetail.Username,
		&userDetail.Name,
		&userDetail.Phone,
		&userDetail.Email,
	)

	if err != nil {
		utils.LogError("Repo", "func GetUserDetail", err)
		return userDetail, err
	}

	return userDetail, err
}

func (db *userConnection) CreateUser(ctx context.Context, tx *sqlx.Tx, u model.User) (string, error) {
	var (
		createdUser string
	)

	query := `
		INSERT INTO 
			users (username, name, password, email, phone) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING email
			;
	`

	err := db.connection.QueryRowxContext(
		ctx,
		query,
		u.Username,
		u.Name,
		u.Password,
		u.Email,
		u.Phone,
	).Scan(
		&createdUser,
	)

	if err != nil {
		utils.LogError("Repo", "func CreateUser", err)
		return createdUser, err
	}

	return createdUser, err
}

func (db *userConnection) FindByEmail(ctx context.Context, email string) (model.UserResponse, error) {
	var (
		userData model.UserResponse
	)

	query := `
		SELECT 
			u.user_id, u.username, u.email, u.password, u.phone
		FROM
			users AS u
		WHERE
			email = $1;
	`

	err := db.connection.QueryRowxContext(
		ctx,
		query,
		email,
	).Scan(
		&userData.User_id,
		&userData.Username,
		&userData.Email,
		&userData.Password,
		&userData.Phone,
	)

	if err != nil {
		utils.LogError("Repo", "func FindByEmail", err)
		return userData, err
	}

	return userData, err
}
