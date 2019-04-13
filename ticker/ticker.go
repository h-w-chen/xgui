package ticker

import "time"

var (
	ticker = time.NewTicker(20 * time.Millisecond)
	value  = 0
)

// GetYOnTicker produces Y on timer
func GetYOnTicker() int {
	<-ticker.C
	y := value
	value++
	return y
}
