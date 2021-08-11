package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/joho/godotenv"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"
	"log"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Tour struct {
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	CoverUrl    string   `json:"cover_url"`
	AudioUrl    string   `json:"audio_url"`
}

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %v", err)
	}

	components, err := common.InitComponents()
	if err != nil {
		panic(err)
	}

	var jlFile bytes.Buffer
	if err = components.Storage.Download("scraping/audiotours/mywowo-full-run-located.jl", &jlFile); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(&jlFile)
	for scanner.Scan() {
		text := scanner.Text()
		var tour Tour
		if err = json.Unmarshal([]byte(text), &tour); err != nil {
			panic(err)
		}

		var user models.User
		handle := tour.Author
		if handle == "" {
			handle = randomdata.SillyName()
			user = models.User{Email: randomdata.Email(), Handle: handle, Password: []byte(randomdata.RandStringRunes(10))}
			if err = user.Insert(context.Background(), components.Db, boil.Infer()); err != nil {
				panic(err)
			}
		} else {
			u, err := models.Users(qm.Where(fmt.Sprintf("%s = ?", models.UserColumns.Handle), handle)).One(context.Background(), components.Db)
			if err != nil {
				panic(err)
			}
			user = *u
		}

		// TODO: rethink asset schema
		audio := models.Asset{MimeType: "audio/mp3", FileName: randomdata.Letters(10)}
		if err = audio.Insert(context.Background(), components.Db, boil.Infer()); err != nil {
			panic(err)
		}

		walk := models.Walk{
			Title:      tour.Title,
			StartPoint: pgeo.Point{X: tour.Location.Lat, Y: tour.Location.Lng},
			AuthorID:   user.ID,
			AudioID:    null.NewInt(audio.ID, true),
		}
		if err = walk.Insert(context.Background(), components.Db, boil.Infer()); err != nil {
			panic(err)
		}
		break
	}

	walks, err := models.Walks().Count(context.Background(), components.Db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d walks in db", walks)
}
