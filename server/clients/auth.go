package clients

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type AuthClient interface {
	VerifyJWT(context.Context, string) (*uuid2.UUID, string, error)
}

type FirebaseAuthClientImpl struct {
	*auth.Client
}

func (f FirebaseAuthClientImpl) VerifyJWT(ctx context.Context, s string) (*uuid2.UUID, string, error) {
	jwt := strings.Replace(s, "Bearer ", "", 1)
	tok, err := f.VerifyIDTokenAndCheckRevoked(ctx, jwt)
	if err != nil {
		return nil, "", status.Error(codes.PermissionDenied, fmt.Sprintf("auth header provided couldn't be verified: %+v", err))
	}

	phonePayload, ok := tok.Claims["phone_number"]
	if !ok {
		// TODO: investigate why there are multiple phone claims payloads!
		otherPhonePayload, otherOk := tok.Claims["phone"]
		if !otherOk {
			return nil, "", status.Error(codes.Internal, fmt.Sprintf("firebase user has no phone claim"))
		}
		phonePayload = otherPhonePayload
	}

	phone, ok := phonePayload.(string)
	if !ok {
		return nil, "", status.Error(codes.Internal, fmt.Sprintf("firebase user has unhandled phone claim (not a string)"))
	}

	uid, err := uuid2.FromString(tok.UID)
	return &uid, phone, err
}

var _ AuthClient = FirebaseAuthClientImpl{}

func NewFirebaseAuthClient(app *firebase.App) *FirebaseAuthClientImpl {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Panicf("could not instantiate firebase auth service: %+v", err)
	}

	return &FirebaseAuthClientImpl{client}
}
