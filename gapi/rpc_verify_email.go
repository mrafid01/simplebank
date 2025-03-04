package gapi

import (
	"context"

	db "github.com/mrafid01/simplebank/db/sqlc"
	"github.com/mrafid01/simplebank/pb"
	"github.com/mrafid01/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	}

	resultTx, err := server.store.VerifyEmailTx(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	resp := &pb.VerifyEmailResponse{
		IsVerified: resultTx.User.IsEmailVerified,
	}
	return resp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
