package api

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/talkiewalkie/talkiewalkie/common"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/models"

	uuid2 "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories"
)

type UserService struct {
}

func (us UserService) CreateUser(ctx context.Context, input *pb.CreateUserInput) (*pb.User, error) {
	components, ok := ctx.Value("components").(*common.Components)
	if !ok {
		return nil, status.Errorf(codes.Internal, "no components in context")
	}

	existing, err := models.Users(models.UserWhere.Handle.EQ(input.Handle), models.UserWhere.FirebaseUID.EQ(null.StringFromPtr(nil)), models.UserWhere.CreatedAt.GT(time.Now().Add(-time.Hour*24*7))).One(components.Ctx, components.Db)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not check for existing users with this handle: %+v", err)
	}
	if existing != nil {
		if _, err := existing.Delete(components.Ctx, components.Db); err != nil {
			return nil, status.Errorf(codes.Internal, "could not delete old user with this handle: %+v", err)
		}
	}

	me := &models.User{
		PhoneNumber: input.PhoneNumber,
		DisplayName: input.DisplayName,
		Handle:      input.Handle,
	}
	if err := me.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %+v", err)
	}

	return repositories.UserToProto(me), nil
}

func (us UserService) Search(ctx context.Context, input *pb.SearchInput) (*pb.SearchOutput, error) {
	components, _, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	users, err := models.Users(
		qm.Where(fmt.Sprintf("%s like '%$1'", models.UserColumns.Handle), input.Prefix),
		qm.Limit(20),
	).All(components.Ctx, components.Db)
	if err != nil {
		return nil, err
	}

	out := []*pb.User{}
	for _, user := range users {
		out = append(out, repositories.UserToProto(user))
	}
	return &pb.SearchOutput{Users: out}, nil
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

	me.DisplayName = strings.TrimSpace(input.DisplayName)
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
				Body:  fmt.Sprintf("%s has joined TalkieWalkie! ❤️", user.DisplayName),
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

var _ pb.UserServiceServer = UserService{}

func NewUserService() UserService {
	return UserService{}
}
