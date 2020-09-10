package main

import (
	"io"
	"os"

	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"

	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales/infrastructure/repository"
	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales/service/managers"
)

func main() {

	numOfTopSpenders := 0
	flag.IntVar(&numOfTopSpenders, "numTopSpenders", 3, "Number of top spenders per month")
	numOfMonths := 6
	flag.IntVar(&numOfMonths, "numMonths", 6, "Number of months")
	var inputFilename string
	flag.StringVar(&inputFilename, "inputFilename", "sample-transactions.csv", "CSV File to read from")
	var outputFilename string
	flag.StringVar(&outputFilename, "outputFilename", "output.csv", "Output filename")
	flag.Parse()

	repos, err := repository.NewCSVLedgerRepository(inputFilename)
	if err != nil {
		log.Error().Err(err).Msg("failed to create ledger repository")
	}

	analysisService := managers.NewAnalysisService(repos)

	report, err := analysisService.TopSpenders(numOfTopSpenders, numOfMonths)
	if err != nil {
		log.Error().Err(err).Msg("failed to perform TopSpenders analysis")
	}

	outputFile, err := os.Create(outputFilename)
	if _, err := io.Copy(outputFile, report.FormattedAsCSV()); err != nil {
		log.Error().Err(err).Msg("failed to write output")
	}
}
