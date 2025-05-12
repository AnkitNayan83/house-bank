package gapi

import (
	"context"
	"errors"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/util"
	"github.com/AnkitNayan83/houseBank/validators"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (res *pb.UpdateUserResponse, err error) {

	authPayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	violations := validateUpdateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "user does not have permission to update other user")
	}

	arg := db.UpdateUserParams{
		Username:        util.NewPgText(req.GetUsername()),
		FullName:        util.NewPgText(req.GetFullName()),
		Email:           util.NewPgText(req.GetEmail()),
		EmailVerifiedAt: util.NewPgTime(req.GetEmailVerifiedAt().AsTime()),
	}

	user, err := server.store.UpdateUser(ctx, arg)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign key violation
				return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot update user: %v", err)
	}

	res = &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.FullName != nil {
		if err := validators.ValidateFullname(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validators.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
