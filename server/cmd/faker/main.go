package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/bxcodec/faker/v3"
	"github.com/friendsofgo/errors"
	"github.com/joho/godotenv"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

var (
	phone = flag.String("phone", "", "normalized phone number, e.g. +33685930995")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *phone == "" {
		fmt.Println("you need to provide a normalized phone number e.g. '-phone +33685930995'")
		return
	}

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %+v", err)
	}

	components := common.InitComponents()

	fbuid, err := components.AuthClient.UserUidByPhoneNumber(ctx, *phone)
	if err != nil {
		log.Panicf("could not fetch firebase user: %+v", err)
	}

	u, err := models.Users(models.UserWhere.FirebaseUID.EQ(null.StringFrom(fbuid))).One(ctx, components.Db)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			u = &models.User{
				DisplayName:        os.Getenv("USER"),
				FirebaseUID:        null.StringFrom(fbuid),
				PhoneNumber:        *phone,
				OnboardingFinished: true,
				Locales:            []string{"fr"},
			}
			if err = u.Insert(ctx, components.Db, boil.Infer()); err != nil {
				log.Panicf("could create user for email: %+v", err)
			}
			if err = common.CreateDefaultConversation(components, context.Background(), u); err != nil {
				log.Panic(err)
			}
		} else {
			log.Panicf("could not find user for email: %+v", err)
		}
	}

	for i := 0; i < 10; i += 1 {
		friends := []*models.User{u}
		for j := 0; j < rand.Intn(6)+1; j += 1 {
			friend := &models.User{
				Handle:             faker.Username(),
				DisplayName:        faker.FirstName(),
				FirebaseUID:        null.String{},
				PhoneNumber:        faker.Phonenumber(),
				OnboardingFinished: true,
				Locales:            []string{"fr"},
			}
			if err = friend.Insert(ctx, components.Db, boil.Infer()); err != nil {
				log.Panicf("could not insert new friend: %+v", err)
			}
			friends = append(friends, friend)
		}

		conv := models.Conversation{Name: null.String{}}
		if err = conv.Insert(ctx, components.Db, boil.Infer()); err != nil {
			log.Panicf("could not insert new conv: %+v", err)
		}

		fmt.Printf("[%s] new conv[%s]", u.DisplayName, conv.UUID.String())

		for _, f := range friends {
			uc := models.UserConversation{UserID: f.ID, ConversationID: conv.ID}
			if err = uc.Insert(ctx, components.Db, boil.Infer()); err != nil {
				log.Panicf("could not link user to conv: %+v", err)
			}
		}

		fmt.Printf(", connected %d friends", len(friends))

		numMsgs := rand.Intn(150) + 1
		fmt.Printf(" - with %d messages\n", numMsgs)
		for j := 0; j < numMsgs; j += 1 {
			text := faker.Paragraph()
			if rand.Int31()%2 == 0 {
				text = faker.Sentence()
			}

			frid := rand.Intn(len(friends))
			authorId := friends[frid].ID

			msg := models.Message{
				Type:           models.MessageTypeText,
				Text:           null.StringFrom(text),
				AuthorID:       null.IntFrom(authorId),
				ConversationID: conv.ID,
			}
			if err = msg.Insert(ctx, components.Db, boil.Infer()); err != nil {
				log.Panicf("could not insert new message: %+v", err)
			}
		}
	}
}
