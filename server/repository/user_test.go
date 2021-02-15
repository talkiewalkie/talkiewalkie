package repository

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setup() PgUserRepository {
	db, err := sqlx.Connect("postgres", "user=theo dbname=talkiewalkie sslmode=disable")

	if err != nil {
		log.Printf("could not connect to db: %v", err)
	}
	return PgUserRepository{Db: db}
}

var (
	userRepo = setup()
)

func TestPgUserRepository_CreateUser(t *testing.T) {
	//mock.ExpectQuery()
	_, err := userRepo.CreateUser("test_email@example.com", []byte("ab1239de"), "secret-xxx")
	if err != nil {
		log.Panicf("floutch %v", err)
	}

}
