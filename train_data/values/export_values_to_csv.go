package train_data

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PlatosCodes/MerolaStation/api"
	db "github.com/PlatosCodes/MerolaStation/db/sqlc"
	"github.com/gin-gonic/gin"
)

func ExportTrainsToCSV(ctx *gin.Context, server *api.Server) error {
	// Fetch all updated trains
	trains, err := server.Store.ListTrains(ctx, db.ListTrainsParams{
		Limit:  10000,
		Offset: 1,
	}) // Assuming you have less than 10,000 trains
	if err != nil {
		return err
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	log.Println(currentTime)
	// Create a CSV file
	new_filepath := fmt.Sprintf("./train_data/values/updated_train%s.csv", currentTime)
	file, err := os.Create(new_filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the headers to the CSV
	err = writer.Write([]string{"ID", "Model Number", "Name", "Value", "Image URL", "Created At", "Version", "Last Edited At"})
	if err != nil {
		return err
	}

	// Write each train record to the CSV
	for _, train := range trains {
		err := writer.Write([]string{
			fmt.Sprintf("%d", train.ID),
			train.ModelNumber,
			train.Name,
			fmt.Sprintf("%d", train.Value),
			train.ImgUrl,
			train.CreatedAt.String(),
			fmt.Sprintf("%d", train.Version),
			train.LastEditedAt.String(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
