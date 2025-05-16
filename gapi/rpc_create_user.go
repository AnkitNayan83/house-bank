package gapi

import (
	"context"
	"errors"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/util"
	"github.com/AnkitNayan83/houseBank/validators"
	"github.com/AnkitNayan83/houseBank/workers"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (res *pb.CreateUserResponse, err error) {

	violations := validateCreateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword()) // using get is safer than using req.Password because it will not panic if the field is not set

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	createUserTxPayload := db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreateUser: func(user db.User) error {
			// Send verification email task to the queue
			taskPayload := &workers.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				// asynq.ProcessIn(10 * time.Second), for testing
				asynq.Queue(workers.QueueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	result, err := server.store.CreateUserTx(ctx, createUserTxPayload)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	res = &pb.CreateUserResponse{
		User: convertUser(result.User),
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validators.ValidateFullname(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := validators.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
