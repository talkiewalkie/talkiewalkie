package common

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errors2 "github.com/friendsofgo/errors"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"math/rand"
	"strings"
)

type Context struct {
	Context    context.Context
	Components *Components
	User       *models.User
}

func AuthInterceptor(c *Components) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("failed to get call metadata")
		}
		jwts := md.Get("Authorization")
		if len(jwts) != 1 {
			return nil, status.Error(codes.PermissionDenied, "missing authorization metadata key")
		}

		tok, err := c.FbAuth.VerifyIDTokenAndCheckRevoked(ctx, strings.Replace(jwts[0], "Bearer ", "", 1))
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("auth header provided couldn't be verified: %+v", err))
		}

		u, err := models.Users(models.UserWhere.FirebaseUID.EQ(null.StringFrom(tok.UID))).One(ctx, c.Db)
		if err != nil && errors2.Cause(err) == sql.ErrNoRows {
			phonePayload, ok := tok.Claims["phone_number"]
			if !ok {
				// TODO: investigate why there are multiple phone claims payloads!
				otherPhonePayload, otherOk := tok.Claims["phone"]
				if !otherOk {
					return nil, status.Error(codes.Internal, fmt.Sprintf("firebase user has no phone claim"))
				}
				phonePayload = otherPhonePayload
			}
			phone, ok := phonePayload.(string)
			if !ok {
				return nil, status.Error(codes.Internal, fmt.Sprintf("firebase user has unhandled phone claim (not a string)"))
			}

			u = &models.User{
				DisplayName:    null.StringFromPtr(nil),
				PhoneNumber:    phone,
				FirebaseUID:    null.NewString(tok.UID, true),
				ProfilePicture: null.NewInt(0, false), // TODO reupload picture
			}
			if err = u.Insert(ctx, c.Db, boil.Infer()); err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("could not create matching db user for new firebase user: %+v", err))
			}
			if err = CreateDefaultConversation(c, ctx, u); err != nil {
				return nil, err
			}

		} else if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to query for user uid: %+v", err))
		}

		newCtx := context.WithValue(ctx, "user", u)
		return newCtx, nil
	}
}

func CreateDefaultConversation(c *Components, ctx context.Context, u *models.User) error {
	// TODO: fetch only once in component init
	twDefaultFriends := []string{"k6WhmQLnpvUCeKuDdpknVzBUu9r1", "YUqVmo08xvXqPZLTYXX7qkvuvGn2"}
	firstFriend := twDefaultFriends[rand.Intn(2)]
	friend, err := models.Users(models.UserWhere.FirebaseUID.EQ(null.StringFrom(firstFriend))).One(ctx, c.Db)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("could not create find default friend: %+v", err))
	}
	firstConv := &models.Conversation{}
	tx, err := c.Db.BeginTx(ctx, nil)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("could not create transaction: %+v", err))
	}
	if err = firstConv.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return status.Error(codes.Internal, fmt.Sprintf("could not create first conversation: %+v", err))
	}
	if err = (&models.UserConversation{ConversationID: firstConv.ID, UserID: friend.ID}).Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return status.Error(codes.Internal, fmt.Sprintf("could not add user to first conversation : %+v", err))
	}
	if err = (&models.UserConversation{ConversationID: firstConv.ID, UserID: u.ID}).Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return status.Error(codes.Internal, fmt.Sprintf("could not add user to first conversation : %+v", err))
	}
	if err = (&models.Message{
		Type:           models.MessageTypeText,
		Text:           null.StringFrom("Hey! This is your open line with the TalkieWalkie team 🤓!"),
		ConversationID: firstConv.ID,
		AuthorID:       null.IntFrom(friend.ID),
	}).Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return status.Error(codes.Internal, fmt.Sprintf("could not insert message in first conversation: %+v", err))
	}
	if err = tx.Commit(); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("could not commit transaction for first conversation: %+v", err))
	}
	return nil
}

func GetUser(ctx context.Context) (*models.User, error) {
	u, ok := ctx.Value("user").(*models.User)
	if !ok {
		return nil, errors.New("No [user] key in context")
	}

	return u, nil
}
