package coco

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type UserService struct {
}

var _ pb.UserServiceServer = UserService{}

func NewUserService() UserService {
	return UserService{}
}

func (us UserService) Get(ctx context.Context, input *pb.UserGetInput) (*pb.User, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	uid, err := uuid2.FromString(input.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse uuid: %+v", err)
	}

	users, err := components.UserStore.ByUuids(uid)
	if err != nil {
		return nil, err
	}

	if ok, err := users[0].HasAccess(me); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user is not in shared conversations: %+v", err)
	}

	return users[0].ToPb(), nil
}

func (us UserService) Me(ctx context.Context, _ *pb.Empty) (*pb.MeUser, error) {
	_, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.MeUser{
		User:         me.ToPb(),
		LanguageUsed: strings.Join(me.Record.Locales, ", "),
	}, nil
}

func (us UserService) Onboarding(ctx context.Context, input *pb.OnboardingInput) (*pb.MeUser, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	me.Record.DisplayName = null.StringFrom(strings.TrimSpace(input.DisplayName))
	me.Record.Locales = input.Locales
	me.Record.OnboardingFinished = true
	if _, err = me.Record.Update(ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %+v", err)
	}

	return &pb.MeUser{
		User:         me.ToPb(),
		LanguageUsed: strings.Join(me.Record.Locales, ", "),
	}, nil
}

func (us UserService) SyncContacts(ctx context.Context, input *pb.SyncContactsInput) (*pb.SyncContactsOutput, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	users, err := components.UserStore.ByPhoneNumbers(input.PhoneNumbers...)
	pbUsers := []*pb.User{}
	for _, user := range users {
		pbUsers = append(pbUsers, user.ToPb())
	}

	if me.Record.BroadcastArrival && len(users) > 0 {
		me.Record.BroadcastArrival = false
		if _, err = me.Record.Update(ctx, components.Db, boil.Infer()); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		messages := []*messaging.Message{}
		for _, user := range users {
			if !user.Record.FirebaseUID.Valid {
				continue
			}
			messages = append(messages, &messaging.Message{
				Topic: user.Record.FirebaseUID.String,
				Data:  map[string]string{"uuid": user.Record.UUID.String()},
				Notification: &messaging.Notification{
					Title: "Good news 🎙!",
					Body:  fmt.Sprintf("%s has joined TalkieWalkie! ❤️", user.DisplayName()),
				},
			})
		}
		res, err := components.FbMssg.SendAll(ctx, messages)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if res.FailureCount > 0 {
			log.Printf("failed to deliver %d messages", res.FailureCount)
		}
	}

	return &pb.SyncContactsOutput{Users: pbUsers}, nil
}
