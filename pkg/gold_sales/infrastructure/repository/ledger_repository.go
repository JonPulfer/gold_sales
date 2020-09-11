package repository

import (
	"github.com/JonPulfer/gold_sales/pkg/gold_sales"
)

// LedgerRepository provides access to stored GoldTransactions.
type LedgerRepository interface {
	FetchAll() ([]gold_sales.GoldPayment, error)
}

type MockLedgerRepository struct {
	ledger MockLedger
}

type MockLedger map[gold_sales.Spender][]gold_sales.GoldPayment

func NewMockLedgerRepository(mockLedger MockLedger) *MockLedgerRepository {
	return &MockLedgerRepository{ledger: mockLedger}
}

func (mlr MockLedgerRepository) FetchAll() ([]gold_sales.GoldPayment, error) {
	goldTransactions := make([]gold_sales.GoldPayment, 0)
	for _, spenderPayments := range mlr.ledger {
		for _, payment := range spenderPayments {
			goldTransactions = append(goldTransactions, payment)
		}
	}
	return goldTransactions, nil
}

type LedgerRepositoryError struct {
	Message string
}

func (lre LedgerRepositoryError) Error() string {
	return lre.Message
}
