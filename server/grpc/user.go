package coco

import (
	"context"
	"errors"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type UserService struct {
	*common.Components
}

var _ pb.UserServiceServer = UserService{}

func NewUserService(c *common.Components) UserService {
	return UserService{Components: c}
}

func (us UserService) Get(ctx context.Context, input *pb.UserGetInput) (*pb.User, error) {
	var u *models.User
	var err error

	switch input.Id.(type) {
	case *pb.UserGetInput_Uuid:
		uid, err := uuid2.FromString(input.GetUuid())
		if err != nil {
			return nil, err
		}

		u, err = models.Users(models.UserWhere.UUID.EQ(uid)).One(ctx, us.Db)
		if err != nil {
			return nil, err
		}
	case *pb.UserGetInput_Handle:
		u, err = models.Users(models.UserWhere.DisplayName.EQ(null.StringFrom(input.GetHandle()))).One(ctx, us.Db)
		if err != nil {
			return nil, err
		}
	case nil:
		return nil, errors.New("input is nil")
	}

	return &pb.User{
		DisplayName: entities.UserDisplayName(u),
		Uuid:        u.UUID.String(),
		Phone:       u.PhoneNumber,
	}, nil
}

func (us UserService) List(input *pb.UserListInput, server pb.UserService_ListServer) error {
	users, err := models.Users(qm.Limit(20), qm.Offset(int(input.Page))).All(server.Context(), us.Db)
	if err != nil {
		return err
	}

	for _, user := range users {
		err = server.Send(&pb.User{
			DisplayName: entities.UserDisplayName(user),
			Uuid:        user.UUID.String(),
			Phone:       user.PhoneNumber,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (us UserService) Me(ctx context.Context, _ *pb.Empty) (*pb.MeUser, error) {
	u, err := common.GetUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &pb.MeUser{
		User: &pb.User{
			DisplayName: entities.UserDisplayName(u),
			Uuid:        u.UUID.String(),
			Phone:       u.PhoneNumber,
		},
		LanguageUsed: "",
	}, nil
}

func (us UserService) Onboarding(ctx context.Context, input *pb.OnboardingInput) (*pb.MeUser, error) {
	u, err := common.GetUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	u.DisplayName = null.StringFrom(strings.TrimSpace(input.DisplayName))
	u.Locales = input.Locales
	u.OnboardingFinished = true
	if _, err = u.Update(ctx, us.Db, boil.Infer()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.MeUser{
		User: &pb.User{
			DisplayName: entities.UserDisplayName(u),
			Uuid:        u.UUID.String(),
			Phone:       u.PhoneNumber,
		},
		LanguageUsed: "",
	}, nil
}

func (us UserService) SyncContacts(ctx context.Context, input *pb.SyncContactsInput) (*pb.SyncContactsOutput, error) {
	_, err := common.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	users, err := models.Users(models.UserWhere.PhoneNumber.IN(input.PhoneNumbers)).All(ctx, us.Db)
	pbUsers := []*pb.User{}
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			DisplayName: entities.UserDisplayName(user),
			Uuid:        user.UUID.String(),
			Phone:       user.PhoneNumber,
		})
	}

	return &pb.SyncContactsOutput{Users: pbUsers}, nil
}
