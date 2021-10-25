package common

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

type Components struct {
	Db          *sqlx.DB
	PgPubSub    *PgPubSub
	EmailClient EmailClient
	JwtAuth     *jwtauth.JWTAuth
	FbAuth      *auth.Client
	FbMssg      *messaging.Client
	Storage     StorageClient
	Audio       *pb.CompressionClient

	CompressImg func(string, int) (string, error)
}

func InitComponents(env string) (*Components, error) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	emailClient := initEmailClient()

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Panicf("could not init the firebase sdk: %+v", err)
	}
	fbAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Panicf("could not instantiate firebase auth service: %+v", err)
	}
	fbMssg, err := app.Messaging(context.Background())
	if err != nil {
		log.Panicf("could not instantiate firebase messaging service: %+v", err)
	}

	storageClient, err := NewGoogleStorageClient(context.Background())
	if err != nil {
		log.Panicf("could not init the storage: %+v", err)
	}

	audioClient, err := NewAudioClient()
	if err != nil {
		// TODO do fail when no audio client
		log.Printf("could not initiate the audio client: %+v", err)
	}

	dbUri := DbUri(
		"talkiewalkie",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		env == "prod")
	db, err := sqlx.Connect("postgres", dbUri)
	if err != nil {
		return nil, fmt.Errorf("could not connect to '%s': %+v", dbUri, err)
	}

	pgPubSub := NewPgPubSub(db, dbUri)

	return &Components{
		Db:          db,
		PgPubSub:    &pgPubSub,
		EmailClient: emailClient,
		JwtAuth:     tokenAuth,
		FbAuth:      fbAuth,
		FbMssg:      fbMssg,
		Storage:     storageClient,
		Audio:       &audioClient,
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
	}, nil
}
