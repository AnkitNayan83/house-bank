package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/AnkitNayan83/houseBank/db/sqlc"
	"github.com/AnkitNayan83/houseBank/pb"
	"github.com/AnkitNayan83/houseBank/util"
	"github.com/AnkitNayan83/houseBank/validators"
	"github.com/AnkitNayan83/houseBank/workers"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (res *pb.LoginUserResponse, err error) {

	violations := validateLoginUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUserByUsername(ctx, req.GetUsername())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}

	err = util.CheckPasswordHash(req.GetPassword(), user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password: %v", err)
	}

	// check if user is verified
	if !user.EmailVerifiedAt.Valid {
		taskPayload := &workers.PayloadSendVerifyEmail{
			Username: user.Username,
		}
		opts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.Queue(workers.QueueueCritical),
		}
		err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot send verification email: %v", err)
		}

		return nil, status.Errorf(codes.Unauthenticated, "user is not verified, verification email sent")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.ACCESS_TOKEN_DURATION)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.REFRESH_TOKEN_DURATION)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	metaData := server.extractMetaData(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metaData.UserAgent,
		ClientID:     metaData.ClientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshTokenPayload.ExpiresAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	res = &pb.LoginUserResponse{
		User:                  convertUser(user),
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  timestamppb.New(accessTokenPayload.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: timestamppb.New(refreshTokenPayload.ExpiresAt),
		SessionId:             session.ID.String(),
	}

	return res, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
