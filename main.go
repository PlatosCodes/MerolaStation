package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	"github.com/PlatosCodes/MerolaStation/api"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/PlatosCodes/MerolaStation/util"
	"github.com/gin-gonic/gin"
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

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	// FOR ADDING TRAIN DATA
	ctx := &gin.Context{}
	loadCSVDataToDB(ctx, server)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func loadCSVDataToDB(ctx *gin.Context, server *api.Server) {
	file, err := os.Open("train_data.csv")
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", "train_data.csv", err.Error())
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','

	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Cannot read '%s': %s\n", "train_data.csv", err.Error())
	}

	for _, line := range lines {
		arg := db.CreateTrainParams{
			ModelNumber: line[1],
			Name:        line[2],
		}

		_, err := server.Store.CreateTrain(ctx, arg)
		if err != nil {
			log.Fatalf("Cannot create train: %s\n", err.Error())
		}
	}
}
