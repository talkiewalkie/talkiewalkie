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

	"github.com/talkiewalkie/talkiewalkie/pb"
)

type Components struct {
	EmailClient EmailClient
	JwtAuth     *jwtauth.JWTAuth
	Storage     StorageClient
	CompressImg func(string, int) (string, error)
	Audio       pb.CompressionClient
}

func InitComponents() Components {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	emailClient := initEmailClient()

	storageClient, err := initStorageClient(context.Background())
	if err != nil {
		log.Panicf("could not init the storage: %+v", err)
	}

	audioClient, err := NewAudioClient()
	if err != nil {
		log.Panicf("could not initiate the audio client: %+v", err)
	}

	return Components{
		EmailClient: emailClient,
		JwtAuth:     tokenAuth,
		Storage:     storageClient,
		Audio:       audioClient,
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
}
