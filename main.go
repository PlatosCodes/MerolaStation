package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	"github.com/PlatosCodes/MerolaStation/api"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/mailer"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// run db migrations
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store, mailer.New(config.Host, config.Port, config.Sender))
	if err != nil {
		log.Fatal("cannot create server")
	}

	// FOR ADDING TRAIN DATA
	ctx := &gin.Context{}
	trainCount, err := server.Store.GetTotalTrainCount(ctx)
	if err != nil {
		log.Fatal("cannot get initial train count", err)
	}
	if trainCount == 0 {
		loadCSVDataToDB(ctx, server)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}

func loadCSVDataToDB(ctx *gin.Context, server *api.Server) {
	file, err := os.Open("./trains/train_data.csv")
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", "./trains/train_data.csv", err.Error())
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','

	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Cannot read '%s': %s\n", "./trains/train_data.csv", err.Error())
	}

	//For train csv with no images
	// for _, line := range lines {
	// 	arg := db.CreateTrainParams{
	// 		ModelNumber: line[1],
	// 		Name:        line[2],
	// 	}

	// 	_, err := server.Store.CreateTrain(ctx, arg)
	// 	if err != nil {
	// 		log.Fatalf("Cannot create train: %s\n", err.Error())
	// 	}
	// }

	//For csv that contains trains with image links
	for _, line := range lines {
		arg := db.CreateImageTrainParams{
			ModelNumber: line[1],
			Name:        line[2],
		}

		if len(line) < 4 {
			log.Printf("Warning: Incomplete data in row: %v\n", line)
			arg.ImgUrl = "" // Setting ImgUrl to empty string
		} else {
			arg.ImgUrl = line[3]
		}

		_, err := server.Store.CreateImageTrain(ctx, arg)
		if err != nil {
			log.Fatalf("Cannot create train: %s\n", err.Error())
		}
	}

}
