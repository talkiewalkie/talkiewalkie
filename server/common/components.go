package common

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories"
)

type Components struct {
	Db *sqlx.DB

	AuthClient      clients.AuthClient
	MessagingClient clients.MessagingClient
	PubSubClient    clients.PubSubClient
	StorageClient   clients.StorageClient
	AudioClient     pb.CompressionClient

	// Context sensitive items
	Ctx context.Context
	repositories.Repositories

	CompressImg func(string, int) (string, error)
}

func (components *Components) ResetEntityStores(ctx context.Context) {
	components.Ctx = ctx
	components.Repositories = repositories.New(ctx, components.Db, components.StorageClient, components.PubSubClient)
}

func InitComponents() *Components {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Panicf("could not init the firebase sdk: %+v", err)
	}

	dbUri := DbUri(
		"talkiewalkie",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		false)
	db, err := sqlx.Connect("postgres", dbUri)
	if err != nil {
		log.Panicf("could not connect to '%s': %+v", dbUri, err)
	}

	models.AddUserHook(boil.BeforeInsertHook, func(ctx context.Context, db boil.ContextExecutor, u *models.User) error {
		u.Handle = slug.Make(u.Handle)
		return nil
	})
	models.AddUserHook(boil.BeforeUpdateHook, func(ctx context.Context, db boil.ContextExecutor, u *models.User) error {
		u.Handle = slug.Make(u.Handle)
		return nil
	})
	models.AddUserHook(boil.BeforeUpsertHook, func(ctx context.Context, db boil.ContextExecutor, u *models.User) error {
		u.Handle = slug.Make(u.Handle)
		return nil
	})

	authClient := clients.NewFirebaseAuthClient(app)
	messagingClient := clients.NewFirebaseMessagingClient(app)
	storageClient := clients.NewGoogleStorageClient(context.Background())
	audioClient := clients.NewAudioClient()
	pubSubClient := clients.NewPgPubSub(db, dbUri)

	components := &Components{
		Db: db,

		AuthClient:      authClient,
		MessagingClient: messagingClient,
		PubSubClient:    pubSubClient,
		StorageClient:   storageClient,
		AudioClient:     audioClient,

		CompressImg: func(path string, width int) (string, error) {
			output := fmt.Sprintf("/tmp/%s.png", strconv.FormatInt(rand.Int63(), 10))

			// https://www.smashingmagazine.com/2015/06/efficient-image-resizing-with-imagemagick/
			cmd := exec.Command(
				"convert", path,
				"-filter", "Triangle",
				"-define", "filter:support=2",
				"-resize", strconv.Itoa(width),
				"-unsharp", "0.25x0.25+8+0.065",
				"-dither", "None",
				"-posterize", "136",
				"-quality", "82",
				"-define", "jpeg:fancy-upsampling=off",
				"-define", "png:compression-filter=5",
				"-define", "png:compression-level=9",
				"-define", "png:compression-strategy=1",
				"-define", "png:exclude-chunk=all",
				"-interlace", "none",
				"-colorspace", "sRGB",
				"-strip", output,
			)
			stdout, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("could not run command: %+v\n%v", err, string(stdout))
			}
			return output, nil
		},
	}

	components.ResetEntityStores(context.Background())
	return components
}
