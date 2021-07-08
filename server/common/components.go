package common

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

type Components struct {
	Db          *sqlx.DB
	EmailClient EmailClient
	JwtAuth     *jwtauth.JWTAuth
	Storage     StorageClient
	Audio       *pb.CompressionClient

	CompressImg func(string, int) (string, error)
}

func InitComponents() (*Components, error) {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	emailClient := initEmailClient()

	storageClient, err := initStorageClient(context.Background())
	if err != nil {
		log.Panicf("could not init the storage: %+v", err)
	}

	audioClient, err := NewAudioClient()
	if err != nil {
		// TODO do fail when no audio client
		log.Printf("could not initiate the audio client: %+v", err)
	}

	dsName := fmt.Sprintf(
		"postgres://%s:%s@%s/talkiewalkie?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
	)
	db, err := sqlx.Connect("postgres", dsName)

	if err != nil {
		return nil, err
	}

	return &Components{
		Db:          db,
		EmailClient: emailClient,
		JwtAuth:     tokenAuth,
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
