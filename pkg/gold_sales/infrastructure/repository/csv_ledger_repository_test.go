package repository

import (
	"testing"

	"github.com/JonPulfer/gold_sales_report/pkg/gold_sales"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	testCases := []struct {
		Name          string
		HeadersIn     []string
		ErrorExpected bool
	}{
		{
			"All headers present",
			requiredHeaders,
			false,
		},
		{
			"One header missing",
			[]string{
				"last_name",
				"email",
				"amount",
				"rate",
				"date",
				"description",
				"to_currency",
				"from_currency",
			},
			true,
		},
		{
			"Two headers missing",
			[]string{
				"email",
				"amount",
				"rate",
				"date",
				"description",
				"to_currency",
				"from_currency",
			},
			true,
		},
		{
			"Three headers missing",
			[]string{
				"email",
				"amount",
				"rate",
				"date",
				"description",
				"from_currency",
			},
			true,
		},
	}

	clr := CSVLedgerRepository{
		filename:      "some.csv",
		file:          nil,
		fieldColIndex: make(map[string]int),
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := clr.parseHeaders(tc.HeadersIn)
			if tc.ErrorExpected {
				assert.NotNil(t, err, "expected error")
			} else {
				assert.Nil(t, err, "unexpected error")
			}
		})
	}
}

func TestParseRow(t *testing.T) {

	testCases := []struct {
		Name            string
		Input           [][]string
		ErrorExpected   bool
		PaymentExpected bool
	}{
		{
			"Correct Payment",
			[][]string{
				[]string{
					"first_name",
					"last_name",
					"email",
					"description",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"CARD SPEND",
					"5311",
					"2629.16",
					"GBP",
					"GGM",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			false,
			true,
		},
		{
			"Not payment",
			[][]string{
				[]string{
					"first_name",
					"last_name",
					"email",
					"description",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"CARD SSPEND",
					"5311",
					"2629.16",
					"GBP",
					"GGM",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			false,
			false,
		},
		{
			"Not gold payment",
			[][]string{
				[]string{
					"first_name",
					"last_name",
					"email",
					"description",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"CARD SPEND",
					"5311",
					"2629.16",
					"GBP",
					"GBP",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			false,
			false,
		},
		{
			"transaction missing field",
			[][]string{
				[]string{
					"first_name",
					"last_name",
					"email",
					"description",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"CARD SPEND",
					"5311",
					"GBP",
					"GBP",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			true,
			false,
		},
		{
			"Correct Payment fields in different order",
			[][]string{
				[]string{
					"description",
					"first_name",
					"last_name",
					"email",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"CARD SPEND",
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"5311",
					"94.1784",
					"GBP",
					"GGM",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			false,
			true,
		},
		{
			"Wrong data in amount field",
			[][]string{
				[]string{
					"description",
					"first_name",
					"last_name",
					"email",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"CARD SPEND",
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"5311",
					"stuff",
					"GBP",
					"GGM",
					"47.0892",
					"22/03/2020 13:28",
				},
			},
			true,
			false,
		},
		{
			"wrong data in rate field",
			[][]string{
				[]string{
					"description",
					"first_name",
					"last_name",
					"email",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"CARD SPEND",
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"5311",
					"94.1784",
					"GBP",
					"GGM",
					"bad",
					"22/03/2020 13:28",
				},
			},
			true,
			false,
		},
		{
			"Wrong date format",
			[][]string{
				[]string{
					"first_name",
					"last_name",
					"email",
					"description",
					"merchant_code",
					"amount",
					"from_currency",
					"to_currency",
					"rate",
					"date",
				},
				[]string{
					"Alayna",
					"Sparks",
					"alayna.sparks@mailinator.com",
					"CARD SPEND",
					"5311",
					"2629.16",
					"GBP",
					"GGM",
					"47.0892",
					"22/03 13:28",
				},
			},
			true,
			false,
		},
	}

	clr := CSVLedgerRepository{
		filename:      "some.csv",
		file:          nil,
		fieldColIndex: make(map[string]int),
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := clr.parseHeaders(tc.Input[0])
			if err != nil {
				t.Logf("failed to parse headers: %s", err.Error())
				t.Fail()
			}
			payment, err := clr.parseRow(tc.Input[1])
			if tc.ErrorExpected {
				assert.NotNil(t, err, "expected error")
			} else {
				assert.Nil(t, err, "unexpected error")
			}
			if tc.PaymentExpected {
				require.NotNil(t, payment, "expected payment")
				assert.Equal(t, gold_sales.GoldSpend, payment.Description,
					"wrong transaction type")
			}
		})
	}
}
