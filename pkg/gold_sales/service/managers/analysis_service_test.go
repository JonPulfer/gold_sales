package managers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales"
	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales/infrastructure/repository"
)

func TestSpenderTotalsByMonth(t *testing.T) {

	testCases := []struct {
		Name           string
		Analysis       *AnalysisService
		ExpectedResult SpenderTotalsByReportMonth
		ExpectError    bool
	}{
		{
			"Multiple spenders in one month",
			analysisServiceForTests(multipleSpendersInOneMonth()),
			validMultipleSpendersInOneMonth(),
			false,
		},
		{
			"Multiple spenders in two months",
			analysisServiceForTests(multipleSpendersInTwoMonths()),
			validMultipleSpendersInTwoMonths(),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			payments, err := tc.Analysis.repository.FetchAll()
			if err != nil {
				t.Logf("problem with mock AnalysisService in test: %s", err.Error())
				t.FailNow()
			}

			spendsBySpender := make(SpendsBySpender)
			for _, payment := range payments {
				if _, ok := spendsBySpender[payment.Spender]; !ok {
					spendsBySpender[payment.Spender] = make([]gold_sales.GoldPayment, 0)
				}
				spendsBySpender[payment.Spender] = append(
					spendsBySpender[payment.Spender], payment)
			}

			result := spenderTotalsByMonth(spendsBySpender)

			compareResultWithExpected(t, result, tc.ExpectedResult)
		})
	}
}

func compareResultWithExpected(t *testing.T, result, expected SpenderTotalsByReportMonth) {

	require.Len(t, result, len(expected))

	for expectedSpendMonth, expectedMonthlySpenders := range expected {
		_, ok := result[expectedSpendMonth]
		require.True(t, ok)

		require.Len(t, result[expectedSpendMonth], len(expectedMonthlySpenders))

		for monthlySpender, spendersMonthlySpend := range expectedMonthlySpenders {
			_, ok := result[expectedSpendMonth][monthlySpender]
			require.True(t, ok)

			assert.Equal(t, result[expectedSpendMonth][monthlySpender].TotalSpend,
				spendersMonthlySpend.TotalSpend, result[expectedSpendMonth][monthlySpender])
		}
	}
}

func analysisServiceForTests(mockLedger repository.MockLedger) *AnalysisService {
	mockRepos := repository.NewMockLedgerRepository(mockLedger)

	return NewAnalysisService(mockRepos)
}

func multipleSpendersInOneMonth() repository.MockLedger {
	spendMonth := firstSpendMonth()
	spendMonthRaw, err := time.Parse("Jan 2006", string(spendMonth))
	if err != nil {
		panic(err.Error())
	}

	spenderOne := spenderOneBuilder()
	spenderTwo := spenderTwoBuilder()
	mockLedger := make(repository.MockLedger)
	mockLedger[spenderOne] = []gold_sales.GoldPayment{
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       200.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         spendMonthRaw,
			GramWeight:   5.0,
		},
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       1000.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         spendMonthRaw,
			GramWeight:   50.0,
		},
	}
	mockLedger[spenderTwo] = []gold_sales.GoldPayment{
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       2.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         spendMonthRaw,
			GramWeight:   2.0 / 20.0,
		},
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       4.0,
			Rate:         20.0,
			ToCurrency:   "USD",
			FromCurrency: "GGM",
			Date:         spendMonthRaw,
			GramWeight:   4.0 / 20.0,
		},
	}

	return mockLedger
}

func multipleSpendersInTwoMonths() repository.MockLedger {
	firstSpendMonth := firstSpendMonth()
	firstSpendMonthRaw, err := time.Parse("Jan 2006", string(firstSpendMonth))
	if err != nil {
		panic(err.Error())
	}
	secondSpendMonth := secondSpendMonth()
	secondSpendMonthRaw, err := time.Parse("Jan 2006", string(secondSpendMonth))
	if err != nil {
		panic(err.Error())
	}

	spenderOne := spenderOneBuilder()
	spenderTwo := spenderTwoBuilder()
	mockLedger := make(repository.MockLedger)
	mockLedger[spenderOne] = []gold_sales.GoldPayment{
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       200.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         firstSpendMonthRaw,
			GramWeight:   5.0,
		},
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       1000.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         firstSpendMonthRaw,
			GramWeight:   50.0,
		},
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       2.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         secondSpendMonthRaw,
			GramWeight:   2.0 / 20.0,
		},
		{
			Spender:      spenderOne,
			Description:  "CARD SPEND",
			Amount:       100.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         secondSpendMonthRaw,
			GramWeight:   100.0 / 20.0,
		},
	}
	mockLedger[spenderTwo] = []gold_sales.GoldPayment{
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       2.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         firstSpendMonthRaw,
			GramWeight:   2.0 / 20.0,
		},
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       4.0,
			Rate:         20.0,
			ToCurrency:   "USD",
			FromCurrency: "GGM",
			Date:         firstSpendMonthRaw,
			GramWeight:   4.0 / 20.0,
		},
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       8.0,
			Rate:         20.0,
			ToCurrency:   "GBP",
			FromCurrency: "GGM",
			Date:         secondSpendMonthRaw,
			GramWeight:   8.0 / 20.0,
		},
		{
			Spender:      spenderTwo,
			Description:  "CARD SPEND",
			Amount:       10.0,
			Rate:         20.0,
			ToCurrency:   "USD",
			FromCurrency: "GGM",
			Date:         secondSpendMonthRaw,
			GramWeight:   10.0 / 20.0,
		},
	}

	return mockLedger
}

func spenderOneBuilder() gold_sales.Spender {
	return gold_sales.Spender{
		FirstName: "Spe",
		LastName:  "nd",
		Email:     "spend@mock.com",
	}
}

func spenderTwoBuilder() gold_sales.Spender {
	return gold_sales.Spender{
		FirstName: "Another",
		LastName:  "Spender",
		Email:     "another_spender@mock.com",
	}
}

func firstSpendMonth() gold_sales.ReportMonth {
	spendMonthRaw, err := time.Parse("Jan 2006", "Jun 2020")
	if err != nil {
		panic(err.Error())
	}
	return gold_sales.ReportMonth(spendMonthRaw.Format("Jan 2006"))
}

func secondSpendMonth() gold_sales.ReportMonth {
	spendMonthRaw, err := time.Parse("Jan 2006", "Jul 2020")
	if err != nil {
		panic(err.Error())
	}
	return gold_sales.ReportMonth(spendMonthRaw.Format("Jan 2006"))
}

func validMultipleSpendersInOneMonth() SpenderTotalsByReportMonth {

	spendMonth := firstSpendMonth()
	spenderOne := spenderOneBuilder()
	spenderTwo := spenderTwoBuilder()
	spenderTotalsInOneMonth := make(SpenderTotalsByReportMonth)
	spenderTotalsInOneMonth[spendMonth] = make(map[gold_sales.Spender]gold_sales.MonthlySpend)
	spenderTotalsInOneMonth[spendMonth][spenderOne] = gold_sales.MonthlySpend{
		Spender:    spenderOne,
		TotalSpend: 55.0,
	}
	spenderTotalsInOneMonth[spendMonth][spenderTwo] = gold_sales.MonthlySpend{
		Spender:    spenderTwo,
		TotalSpend: 0.30000000000000004,
	}

	return spenderTotalsInOneMonth
}

func validMultipleSpendersInTwoMonths() SpenderTotalsByReportMonth {

	firstSpendMonth := firstSpendMonth()
	spenderOne := spenderOneBuilder()
	spenderTwo := spenderTwoBuilder()
	spenderTotalsInOneMonth := make(SpenderTotalsByReportMonth)
	spenderTotalsInOneMonth[firstSpendMonth] = make(map[gold_sales.Spender]gold_sales.MonthlySpend)
	spenderTotalsInOneMonth[firstSpendMonth][spenderOne] = gold_sales.MonthlySpend{
		Spender:    spenderOne,
		TotalSpend: 55.0,
	}
	spenderTotalsInOneMonth[firstSpendMonth][spenderTwo] = gold_sales.MonthlySpend{
		Spender:    spenderTwo,
		TotalSpend: 0.30000000000000004,
	}
	secondSpendMonth := secondSpendMonth()
	spenderTotalsInOneMonth[secondSpendMonth] = make(map[gold_sales.Spender]gold_sales.MonthlySpend)
	spenderTotalsInOneMonth[secondSpendMonth][spenderOne] = gold_sales.MonthlySpend{
		Spender:    spenderOne,
		TotalSpend: 5.1,
	}
	spenderTotalsInOneMonth[secondSpendMonth][spenderTwo] = gold_sales.MonthlySpend{
		Spender:    spenderTwo,
		TotalSpend: 0.9,
	}

	return spenderTotalsInOneMonth
}
