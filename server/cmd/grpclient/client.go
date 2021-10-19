package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
)

var (
	convUuid = flag.String("convUuid", "", "")
	audio    = flag.String("audio", "", "")
	token    = flag.String("token", "", "")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	client := pb.NewMessageServiceClient(conn)
	audioContent, err := ioutil.ReadFile(*audio)
	if err != nil {
		log.Panic(err)
	}

	authHeader := metadata.New(map[string]string{"Authorization": fmt.Sprintf("Bearer %s", *token)})
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
