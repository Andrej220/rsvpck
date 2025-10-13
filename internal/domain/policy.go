package domain

type ExecutionPolicy int

const (
	PolicyOptimized ExecutionPolicy = iota
	PolicyExhaustive
)

var policyNames = [...]string{
	"Optimized",
	"Exhaustive",
}

func (e ExecutionPolicy) String() string {
	if int(e) < len(policyNames) {
		return policyNames[e]
	}
	return "Unknown"
}
