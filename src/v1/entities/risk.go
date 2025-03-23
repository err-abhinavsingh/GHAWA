package entities

type RiskState string

const (
	Open          RiskState = "open"
	Closed        RiskState = "closed"
	Accepted      RiskState = "accepted"
	Investigating RiskState = "investigating"
)

var RiskStates = []RiskState{Open, Closed, Accepted, Investigating}

var ValidRisk = func(state RiskState) bool {

	var found bool
	found = false
	for _, riskState := range RiskStates {
		if riskState == state {
			found = true
			break
		}
	}
	return found
}

type Risk struct {
	Uuid        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       RiskState `json:"state"`
}
