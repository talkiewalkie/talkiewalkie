package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	convUuid = flag.String("convUuid", "", "")
	audio    = flag.String("audio", "", "")

	host  = flag.String("host", "localhost:8080", "")
	fbUid = flag.String("asUser", "k6WhmQLnpvUCeKuDdpknVzBUu9r1", "firebase uid")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %+v", err)
	}

	components, err := common.InitComponents()
	if err != nil {
		log.Panicf("could not initiate components: %+v", err)
	}

	ctok, err := components.FbAuth.CustomToken(context.Background(), *fbUid)
	if err != nil {
		log.Panic(err)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd cmd/grpclient/tokenfetcher && ./getIdToken.js '%s'", ctok))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic(string(output))
	}
	token := strings.TrimSpace(string(output))

	secOpt := grpc.WithInsecure()
	conn, err := grpc.Dial(*host, secOpt, grpc.WithConnectParams(grpc.ConnectParams{
		MinConnectTimeout: 3 * time.Second,
	}))
	if err != nil {
		log.Panic(err)
	}

	client := pb.NewMessageServiceClient(conn)
	audioContent, err := ioutil.ReadFile(*audio)
	if err != nil {
		log.Panic(err)
	}

	authHeader := metadata.New(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)})
	ctx := metadata.NewOutgoingContext(context.Background(), authHeader)
	if _, err = client.Send(ctx, &pb.MessageSendInput{
		Recipients: &pb.MessageSendInput_ConvUuid{ConvUuid: *convUuid},
		Content: &pb.MessageSendInput_VoiceMessage{VoiceMessage: &pb.VoiceMessage{
			RawContent: audioContent,
			SiriTranscript: &pb.AlignedTranscript{Items: []*pb.TranscriptItem{
				{
					Word:            "Hello",
					OffsetMs:        413,
					DurationMs:      200,
					SubstringOffset: 0,
				},
				{
					Word:            "Alex",
					OffsetMs:        700,
					DurationMs:      200,
					SubstringOffset: 5,
				}}, Rendered: "Hello Alex!"},
		}},
	}); err != nil {
		log.Panicf("failed: %+v", err)
	}
	log.Printf("success!")
}
