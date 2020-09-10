package gold_sales

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// MonthlyTopSpendersAnalysisReport that ranks the top spenders by month.
type MonthlyTopSpendersAnalysisReport struct {
	monthlySpenders map[ReportMonth]MonthlySpenders
	months          OrderedReportMonths
	numOfMonths     int
}

func NewMonthlyTopSpendersAnalysisReport(numOfMonths int) *MonthlyTopSpendersAnalysisReport {
	return &MonthlyTopSpendersAnalysisReport{
		monthlySpenders: make(map[ReportMonth]MonthlySpenders),
		numOfMonths:     numOfMonths,
		months:          make(OrderedReportMonths, 0),
	}
}

func (mtsar *MonthlyTopSpendersAnalysisReport) AddMonth(
	monthToAdd ReportMonth,
	monthlySpenders MonthlySpenders,
) error {
	if _, ok := mtsar.monthlySpenders[monthToAdd]; ok {
		return errors.New("month already in report")
	}
	mtsar.monthlySpenders[monthToAdd] = monthlySpenders
	mtsar.months = append(mtsar.months, monthToAdd)
	return nil
}

func (mtsar MonthlyTopSpendersAnalysisReport) String() string {
	line := ""
	for spendMonth, topSpenders := range mtsar.monthlySpenders {
		line = line + fmt.Sprintf("Month: %s\nSpenders: %v\n",
			spendMonth, topSpenders)
	}
	return line
}

// FormattedAsCSV in a buffer ready to be copied to an io.Writer.
func (mtsar *MonthlyTopSpendersAnalysisReport) FormattedAsCSV() *bytes.Buffer {
	var buf bytes.Buffer
	sort.Sort(mtsar.months)

	for i := 0; i < mtsar.numOfMonths && i < len(mtsar.months); i++ {
		for _, monthlySpend := range mtsar.monthlySpenders[mtsar.months[i]] {
			line := fmt.Sprintf("%s,%s,%s,%.2f,\n",
				mtsar.months[i],
				monthlySpend.Spender.FirstName,
				monthlySpend.Spender.LastName,
				monthlySpend.TotalSpend,
			)
			buf.WriteString(line)
		}
	}

	return &buf
}

// MonthlySpenders lists the monthly spenders and can be sorted by spend.
type MonthlySpenders []MonthlySpend

func (ms MonthlySpenders) Len() int {
	return len(ms)
}
func (ms MonthlySpenders) Less(i, j int) bool {
	return ms[i].TotalSpend > ms[j].TotalSpend
}
func (ms MonthlySpenders) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

// MonthlySpend for a particular Spender.
type MonthlySpend struct {
	Spender    Spender    `json:"spender"`
	TotalSpend TotalSpend `json:"totalSpend"`
}

// TotalSpend formatted to meet the business requirements.
type TotalSpend float64

func (ts TotalSpend) String() string {
	return fmt.Sprintf("%.2f", ts)
}

// ReportMonth in the format of `MMM YYYY`.
type ReportMonth string

// ParseReportMonth from a time.Time.
func ParseReportMonth(from time.Time) ReportMonth {
	return ReportMonth(from.Format("Jan 2006"))
}

type OrderedReportMonths []ReportMonth

func (orm OrderedReportMonths) Len() int {
	return len(orm)
}
func (orm OrderedReportMonths) Less(i, j int) bool {
	left, _ := time.Parse("Jan 2006", string(orm[i]))
	right, _ := time.Parse("Jan 2006", string(orm[j]))
	return left.UnixNano() > right.UnixNano()
}
func (orm OrderedReportMonths) Swap(i, j int) {
	orm[i], orm[j] = orm[j], orm[i]
}
