package api

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories"
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

	users, err := components.UserRepository.ByUuids(uid)
	if err != nil {
		return nil, err
	}

	if ok, err := components.UserHasAccess(me, users...); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user is not in shared conversations: %+v", err)
	}

	return repositories.UserToProto(users[0]), nil
}

func (us UserService) Me(ctx context.Context, _ *pb.Empty) (*pb.MeUser, error) {
	_, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.MeUser{
		User:         repositories.UserToProto(me),
		LanguageUsed: strings.Join(me.Locales, ", "),
	}, nil
}

func (us UserService) Onboarding(ctx context.Context, input *pb.OnboardingInput) (*pb.MeUser, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	me.DisplayName = null.StringFrom(strings.TrimSpace(input.DisplayName))
	me.Locales = input.Locales
	me.OnboardingFinished = true
	if _, err = me.Update(ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %+v", err)
	}

	return &pb.MeUser{
		User:         repositories.UserToProto(me),
		LanguageUsed: strings.Join(me.Locales, ", "),
	}, nil
}

func (us UserService) SyncContacts(ctx context.Context, input *pb.SyncContactsInput) (*pb.SyncContactsOutput, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	users, err := components.UserRepository.ByPhoneNumbers(input.PhoneNumbers...)
	pbUsers := []*pb.User{}
	for _, user := range users {
		pbUsers = append(pbUsers, repositories.UserToProto(user))
	}

	if me.BroadcastArrival && len(users) > 0 {
		me.BroadcastArrival = false
		if _, err = me.Update(ctx, components.Db, boil.Infer()); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		messages := []clients.MessageInput{}
		for _, user := range users {
			if !user.FirebaseUID.Valid {
				continue
			}
			messages = append(messages, clients.MessageInput{
				Topic: user.FirebaseUID.String,
				Data:  map[string]string{"uuid": user.UUID.String()},
				Title: "Good news 🎙!",
				Body:  fmt.Sprintf("%s has joined TalkieWalkie! ❤️", repositories.UserDisplayName(user)),
			})
		}
		res, err := components.MessagingClient.SendAll(ctx, messages)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if res.FailureCount > 0 {
			log.Printf("failed to deliver %d messages", res.FailureCount)
		}
	}

	return &pb.SyncContactsOutput{Users: pbUsers}, nil
}
