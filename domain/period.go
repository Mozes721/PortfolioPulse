package domain

// Period represents a valid reporting window.
type Period string

const (
	Period1D  Period = "1d"
	Period1W  Period = "1w"
	Period1M  Period = "1m"
	Period3M  Period = "3m"
	Period1Y  Period = "1y"
	PeriodMax Period = "max"
)

var validPeriods = map[Period]struct{}{
	Period1D:  {},
	Period1W:  {},
	Period1M:  {},
	Period3M:  {},
	Period1Y:  {},
	PeriodMax: {},
}

// Validate returns ErrInvalidPeriod if p is not a recognised token.
func (p Period) Validate() error {
	if _, ok := validPeriods[p]; !ok {
		return ErrInvalidPeriod
	}
	return nil
}
