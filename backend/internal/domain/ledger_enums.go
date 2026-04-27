package domain

type AccountType int
type Direction int
type TxnType int
type Status int
type Currency int

const (
	CHECKING AccountType = iota
	ASSET
)
const (
	DEBIT Direction = iota
	CREDIT
)
const (
	DEPOSIT TxnType = iota
	PAYOUT
)
const (
	PENDING Status = iota
	PROCESSING
	SUCCESS
	FAILED
)
const (
	USD Currency = iota
	GBP
	NGN
	KES
)

func (a AccountType) String() string {
	options := map[AccountType]string{
		CHECKING: "CHECKING",
		ASSET:    "ASSET",
	}

	return options[a]
}

func (d Direction) String() string {
	options := map[Direction]string{
		DEBIT:  "DEBIT",
		CREDIT: "CREDIT",
	}

	return options[d]
}

func (t TxnType) String() string {
	options := map[TxnType]string{
		DEPOSIT: "DEPOSIT",
		PAYOUT:  "PAYOUT",
	}

	return options[t]
}

func (s Status) String() string {
	options := map[Status]string{
		PENDING:    "PENDING",
		PROCESSING: "PROCESSING",
		SUCCESS:    "SUCCESS",
		FAILED:     "FAILED",
	}

	return options[s]
}

func (c Currency) String() string {
	options := map[Currency]string{
		USD: "USD",
		GBP: "GBP",
		NGN: "NGN",
		KES: "KES",
	}

	return options[c]
}
func GetCurrency(s string) (Currency, bool) {
	options := map[string]Currency{
		"USD": USD,
		"GBP": GBP,
		"NGN": NGN,
		"KES": KES,
	}
	currency, ok := options[s]

	return currency, ok
}
