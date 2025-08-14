package presentation

import "time"

// nowFunc is a function that returns the current time.
// Modelled as a global variable to allow for easy mocking in tests.
var nowFunc func() time.Time

func resetNowFunc() {
	nowFunc = time.Now
}

func init() {
	resetNowFunc()
}
