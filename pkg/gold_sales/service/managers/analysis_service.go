package managers

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/JonPulfer/gold_sales/pkg/gold_sales"
	"github.com/JonPulfer/gold_sales/pkg/gold_sales/infrastructure/repository"
)

// AnalysisService performs the high level operations that the business
// requires.
type AnalysisService struct {
	repository repository.LedgerRepository
}

func NewAnalysisService(repository repository.LedgerRepository) *AnalysisService {
	return &AnalysisService{repository: repository}
}

// TopSpenders is a report of the top 3 spenders each month for the last 6
// months.
func (ts AnalysisService) TopSpenders(
	numberSpenders int,
	numberMonths int,
) (
	*gold_sales.MonthlyTopSpendersAnalysisReport,
	error,
) {

	payments, err := ts.repository.FetchAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get payments from repository")
	}

	groupedSpends, err := groupTotalSpendsByMonth(payments)
	if err != nil {
		return nil, errors.Wrap(err, "failed to group the monthly spends")
	}

	monthlyTopSpenders, err := monthlySpenders(groupedSpends, numberSpenders, numberMonths)
	if err != nil {
		return nil, err
	}

	return monthlyTopSpenders, nil
}

func monthlySpenders(
	groupedSpends map[gold_sales.ReportMonth]gold_sales.MonthlySpenders,
	numberSpenders int,
	numOfMonths int,
) (
	*gold_sales.MonthlyTopSpendersAnalysisReport,
	error,
) {

	report := gold_sales.NewMonthlyTopSpendersAnalysisReport(numOfMonths)
	for spendMonth, spenders := range groupedSpends {
		sort.Sort(spenders)

		topMonthSpenders := make(gold_sales.MonthlySpenders, 0)
		for i := 0; i < numberSpenders && i < len(spenders); i++ {
			topMonthSpenders = append(topMonthSpenders, spenders[i])
		}

		err := report.AddMonth(spendMonth, topMonthSpenders)
		if err != nil {
			return gold_sales.NewMonthlyTopSpendersAnalysisReport(numOfMonths), err
		}
	}

	return report, nil
}

// groupTotalSpendsByMonth collates the payments into month spends for each
// spender.
func groupTotalSpendsByMonth(
	payments []gold_sales.GoldPayment,
) (
	map[gold_sales.ReportMonth]gold_sales.MonthlySpenders,
	error,
) {

	spendsBySpender := make(SpendsBySpender)
	for _, payment := range payments {
		if _, ok := spendsBySpender[payment.Spender]; !ok {
			spendsBySpender[payment.Spender] = make([]gold_sales.GoldPayment, 0)
		}
		spendsBySpender[payment.Spender] = append(
			spendsBySpender[payment.Spender], payment)
	}

	spenderTotalsByMonth := spenderTotalsByMonth(spendsBySpender)

	monthlySpends := make(map[gold_sales.ReportMonth]gold_sales.MonthlySpenders)

	for spenderMonth, spenders := range spenderTotalsByMonth {
		for _, spenderMonthlySpend := range spenders {
			if _, ok := monthlySpends[spenderMonth]; !ok {
				monthlySpends[spenderMonth] = make(gold_sales.MonthlySpenders, 0)
			}
			monthlySpends[spenderMonth] = append(monthlySpends[spenderMonth],
				spenderMonthlySpend)
		}
	}

	return monthlySpends, nil
}

// spenderTotalsByMonth collates the monthly spends for each Spender and builds
// Map keyed by ReportMonth to provide easy access to the monthly data.
func spenderTotalsByMonth(
	spendsBySpender SpendsBySpender,
) SpenderTotalsByReportMonth {

	spenderTotalsByMonth := make(SpenderTotalsByReportMonth)
	for spender, spends := range spendsBySpender {
		monthlySpends := make(map[gold_sales.ReportMonth][]gold_sales.GoldPayment)

		for _, spend := range spends {
			spendMonth := gold_sales.ParseReportMonth(spend.Date)
			if _, ok := monthlySpends[spendMonth]; !ok {
				monthlySpends[spendMonth] = make([]gold_sales.GoldPayment, 0)
			}
			monthlySpends[spendMonth] = append(monthlySpends[spendMonth], spend)
		}

		for spendMonth, monthSpends := range monthlySpends {
			var totalWeight float64
			for _, spend := range monthSpends {
				totalWeight = totalWeight + spend.GramWeight
			}
			monthlySpend := gold_sales.MonthlySpend{
				Spender:    spender,
				TotalSpend: gold_sales.TotalSpend(totalWeight),
			}
			spenderTotalsByMonth.Insert(spendMonth, monthlySpend)
		}
	}
	return spenderTotalsByMonth
}

// SpendsBySpender indexes all the spends by the Spender.
type SpendsBySpender map[gold_sales.Spender][]gold_sales.GoldPayment

// SpenderTotalsByReportMonth indexes the Spender totals by ReportMonth.
type SpenderTotalsByReportMonth map[gold_sales.ReportMonth]map[gold_sales.Spender]gold_sales.MonthlySpend

// Insert a new MonthlySpend in the ReportMonth for the Spender in the MonthlySpend.
func (stbrm SpenderTotalsByReportMonth) Insert(
	spendMonth gold_sales.ReportMonth, monthlySpend gold_sales.MonthlySpend) {

	if _, ok := stbrm[spendMonth]; !ok {
		stbrm[spendMonth] = make(map[gold_sales.Spender]gold_sales.MonthlySpend)
	}

	if _, ok := stbrm[spendMonth][monthlySpend.Spender]; !ok {
		stbrm[spendMonth][monthlySpend.Spender] = monthlySpend
	}
}
