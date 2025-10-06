package domain

type ExecutionPolicy int

const (
	PlicyOptimized 		ExecutionPolicy = iota
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