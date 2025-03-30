package util

const (
	USD = "USD"
	INR = "INR"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, INR:
		return true
	default:
		return false
	}
}
