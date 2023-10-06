package currency

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	SEK = "SEK"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, SEK:
		return true
	}
	return false
}
