package operator

func getErrorType(id int) string {
	switch id {
	case 1:
		return "Signature invalid"
	case 2:
		return "Maker signature invalid"
	case 3:
		return "Taker signature invalid"
	case 4:
		return "Orders should have opposite side"
	case 5:
		return "Pricepoints do no match"
	case 6:
		return "Trades already completed or cancelled"
	case 7:
		return "Trade amount is too large"
	case 8:
		return "Rounding error is too large"
	default:
		return "Unknown error"
	}
}
