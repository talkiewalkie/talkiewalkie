package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"strings"
)

var (
	convUuid = flag.String("convUuid", "", "")
	audio    = flag.String("audio", "", "")
	token    = flag.String("token", "", "")
	host     = flag.String("host", "localhost:8080", "")
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	var secOpt grpc.DialOption
	if strings.HasPrefix(*host, "localhost") {
		secOpt = grpc.WithInsecure()
	} else {
		conn, err := tls.Dial("tcp", *host, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			panic(err)
		}

		certs := conn.ConnectionState().PeerCertificates
		_ = conn.Close()

		secOpt = grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&tls.Certificate{
			Certificate: [][]byte{certs[0].Raw},
		}))
	}

	conn, err := grpc.Dial(*host, secOpt)
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
