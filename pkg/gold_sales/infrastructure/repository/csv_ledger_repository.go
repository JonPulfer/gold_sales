package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales"
)

// CSVLedgerRepository uses a CSV file as a LedgerRepository for Gold Payments.
type CSVLedgerRepository struct {
	filename      string
	file          *os.File
	fieldColIndex map[string]int
}

// NewCSVLedgerRepository opens the provided CSV file as a LedgerRepository.
func NewCSVLedgerRepository(filename string) (*CSVLedgerRepository, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	colIndex := make(map[string]int)
	return &CSVLedgerRepository{
		filename:      filename,
		file:          file,
		fieldColIndex: colIndex}, nil
}

func (clr CSVLedgerRepository) FetchAll() ([]gold_sales.GoldPayment, error) {
	goldPayments := make([]gold_sales.GoldPayment, 0)

	rdr := csv.NewReader(clr.file)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, err
	}

	if err := clr.parseHeaders(rows[0]); err != nil {
		return nil, err
	}

	for _, row := range rows[1:] {
		payment, err := clr.parseRow(row)
		if err != nil {
			return nil, err
		}
		if payment != nil {
			goldPayments = append(goldPayments, *payment)
		}
	}

	return goldPayments, nil
}

// requiredHeaders we need to find in the CSV file to be able to extract the
// payment information.
var requiredHeaders = []string{
	"first_name",
	"last_name",
	"email",
	"amount",
	"rate",
	"date",
	"description",
	"to_currency",
	"from_currency",
}

// cleanNonPrintable characters that may sneak in to the headers.
var cleanNonPrintable = regexp.MustCompile("[^a-z_A-Z0-9]+")

// parseHeaders validates we have all the expected headers in the CSV and records
// the column index for each header.
func (clr *CSVLedgerRepository) parseHeaders(headers []string) error {
	requiredHeaderFound := make(map[string]bool)
	for _, requiredHeader := range requiredHeaders {
		requiredHeaderFound[requiredHeader] = false
	}
	for colIdx, header := range headers {
		header = cleanNonPrintable.ReplaceAllString(header, "")
		if _, ok := requiredHeaderFound[header]; ok {
			requiredHeaderFound[header] = true
		}
		clr.fieldColIndex[header] = colIdx
	}

	fieldsNotFound := make([]string, 0)
	for requiredHeader, found := range requiredHeaderFound {
		if !found {
			fieldsNotFound = append(fieldsNotFound, requiredHeader)
		}
	}
	if len(fieldsNotFound) > 0 {
		return LedgerRepositoryError{
			Message: "failed to find the following fields in the CSV: " +
				fmt.Sprintf("%s", fieldsNotFound),
		}
	}

	return nil
}

func (clr CSVLedgerRepository) parseRow(row []string) (*gold_sales.GoldPayment, error) {
	if len(row) != len(clr.fieldColIndex) {
		return nil, LedgerRepositoryError{
			Message: "failed to parse row, unexpected number of fields"}
	}

	amount, err := strconv.ParseFloat(row[clr.fieldColIndex["amount"]], 64)
	if err != nil {
		return nil, LedgerRepositoryError{
			Message: "failed to parse amount: " +
				row[clr.fieldColIndex["amount"]],
		}
	}

	rate, err := strconv.ParseFloat(row[clr.fieldColIndex["rate"]], 64)
	if err != nil {
		return nil, LedgerRepositoryError{
			Message: "failed to parse rate: " +
				row[clr.fieldColIndex["rate"]],
		}
	}

	date, err := time.Parse("02/01/2006 15:04", row[clr.fieldColIndex["date"]])
	if err != nil {
		return nil, LedgerRepositoryError{
			Message: "failed to parse date: " +
				row[clr.fieldColIndex["date"]],
		}
	}

	transaction := gold_sales.GoldPayment{
		Spender: gold_sales.Spender{
			FirstName: row[clr.fieldColIndex["first_name"]],
			LastName:  row[clr.fieldColIndex["last_name"]],
			Email:     row[clr.fieldColIndex["email"]],
		},
		Description:  row[clr.fieldColIndex["description"]],
		Amount:       amount,
		Rate:         rate,
		FromCurrency: row[clr.fieldColIndex["from_currency"]],
		ToCurrency:   row[clr.fieldColIndex["to_currency"]],
		Date:         date,
		GramWeight:   amount / rate,
	}

	if transaction.Description == gold_sales.GoldSpend &&
		transaction.ToCurrency == gold_sales.GoldCurrencyCode {
		return &transaction, nil
	}

	return nil, nil
}
