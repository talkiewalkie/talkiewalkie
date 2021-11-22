package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/pb"
)

var (
	convUuid = flag.String("convUuid", "", "")
	audio    = flag.String("audio", "", "")
	msg      = flag.String("msg", "", "")

	host  = flag.String("host", "localhost:8080", "")
	fbUid = flag.String("asUser", "k6WhmQLnpvUCeKuDdpknVzBUu9r1", "firebase uid")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %+v", err)
	}

	components := common.InitComponents()

	ctok, err := components.AuthClient.CustomToken(context.Background(), *fbUid)
	if err != nil {
		log.Panic(err)
	}

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

	client := pb.NewEventServiceClient(conn)
	var content pb.MessageSendInputContentOneOf
	if *audio != "" {
		content = voiceMessage(*audio)
	} else if *msg != "" {
		content = &pb.MessageSendInput_TextMessage{TextMessage: &pb.TextMessage{Content: *msg}}
	} else {
		log.Panic("need -msg or -audio to be set")
	}

	authHeader := metadata.New(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)})
	ctx := metadata.NewOutgoingContext(context.Background(), authHeader)
	if _, err = client.Sync(
		ctx,
		&pb.UpSync{
			Events: []*pb.Event{{
				Content: &pb.Event_SentNewMessage_{
					SentNewMessage: &pb.Event_SentNewMessage{
						Message:      &pb.MessageSendInput{Content: content},
						Conversation: &pb.Event_SentNewMessage_ConvUuid{ConvUuid: *convUuid}},
				}}},
		},
	); err != nil {
		log.Panicf("failed: %+v", err)
	}
	log.Printf("success!")
}

func voiceMessage(p string) *pb.MessageSendInput_VoiceMessage {
	audioContent, err := ioutil.ReadFile(*audio)
	if err != nil {
		log.Panic(err)
	}

	return &pb.MessageSendInput_VoiceMessage{VoiceMessage: &pb.VoiceMessage{
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
	}}
}
