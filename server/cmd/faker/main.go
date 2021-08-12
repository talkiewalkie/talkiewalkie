package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"
	"log"
	"strings"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type TourFile struct {
	Url      string `json:"url"`
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}
type Tour struct {
	Location    Location   `json:"location"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	CoverUrl    string     `json:"cover_url"`
	AudioUrl    string     `json:"audio_url"`
	TourUrl     string     `json:"tour_url"`
	Files       []TourFile `json:"files"`
}

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %+v", err)
	}

	components, err := common.InitComponents()
	if err != nil {
		log.Panicf("could not initiate components: %+v", err)
	}

	var f bytes.Buffer
	if err = components.Storage.Download("scraping/audiotours/mywowo-full-run-located.jl", &f); err != nil {
		log.Panicf("could not download asset file: %+v", err)
	}

	mywowo := models.User{Handle: "mywowo"}
	if err = mywowo.Insert(ctx, components.Db, boil.Infer()); err != nil {
		log.Panicf("-!- could not insert new user with handle '%s': %+v", "mywowo", err)
	}

	scanner := bufio.NewScanner(&f)
	cnt := 0
	for scanner.Scan() {
		text := scanner.Text()
		var tour Tour
		if err = json.Unmarshal([]byte(text), &tour); err != nil {
			log.Printf("-!- could not deserialize json line: %+v", err)
			continue
		}

		var user models.User
		handle := tour.Author
		if handle == "" {
			handle = mywowo.Handle
			user = mywowo
		} else {
			exists, err := models.Users(qm.Where(fmt.Sprintf("%s = ?", models.UserColumns.Handle), handle)).Exists(ctx, components.Db)
			if err != nil {
				log.Printf("-!- could not look for user with handle '%s': %+v", handle, err)
				continue
			}
			if !exists {
				user = models.User{Handle: handle}
				if err = user.Insert(ctx, components.Db, boil.Infer()); err != nil {
					log.Printf("-!- could not insert new user with handle '%s': %+v", handle, err)
					continue
				}
			} else {
				u, err := models.Users(qm.Where(fmt.Sprintf("%s = ?", models.UserColumns.Handle), handle)).One(ctx, components.Db)
				if err != nil {
					log.Printf("-!- could not retrieve user '%s': %+v", handle, err)
					continue
				}
				user = *u
			}
		}

		var audioFile, coverFile TourFile
		for _, file := range tour.Files {
			if strings.HasSuffix(file.Url, ".mp3") {
				audioFile = file
			} else {
				coverFile = file
			}
		}

		audio := models.Asset{
			MimeType: "audio/mp3",
			FileName: fmt.Sprintf("%s.mp3", audioFile.Checksum),
			Bucket:   null.NewString("talkiewalkie-dev", true),
			BlobName: null.NewString(fmt.Sprintf("scraping/audiotours/%s", audioFile.Path), true),
		}
		if err = audio.Insert(ctx, components.Db, boil.Infer()); err != nil {
			log.Printf("-!- could not insert audio track as asset: %+v", err)
			continue
		}

		cover := models.Asset{
			MimeType: "image/jpg",
			FileName: fmt.Sprintf("%s.mp3", audioFile.Checksum),
			Bucket:   null.NewString("talkiewalkie-dev", true),
			BlobName: null.NewString(fmt.Sprintf("scraping/audiotours/%s", coverFile.Path), true),
		}
		if err = cover.Insert(ctx, components.Db, boil.Infer()); err != nil {
			log.Printf("-!- could not insert cover image as asset: %+v", err)
			continue
		}

		walk := models.Walk{
			Title:       tour.Title,
			Description: null.NewString(tour.Description, true),
			StartPoint:  pgeo.Point{X: tour.Location.Lat, Y: tour.Location.Lng},
			AuthorID:    user.ID,
			AudioID:     null.NewInt(audio.ID, true),
			CoverID:     null.NewInt(cover.ID, true),
		}
		if err = walk.Insert(ctx, components.Db, boil.Infer()); err != nil {
			log.Printf("-!- could not insert walk: %+v", err)
			continue
		}

		cnt += 1
		if cnt%50 == 0 {
			log.Printf("%d walks inserted", cnt)
		}
	}

	walks, err := models.Walks().Count(ctx, components.Db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d walks in db", walks)
}
