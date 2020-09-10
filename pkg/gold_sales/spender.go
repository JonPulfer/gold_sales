package gold_sales

// Spender identifies the actor a gold transaction relates to.
type Spender struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (s Spender) String() string {
	return s.Email
}
