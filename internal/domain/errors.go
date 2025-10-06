package domain


type ErrorCode int

const (
	ErrorCodeDNSUnresolvable ErrorCode = iota
	ErrorCodeTCPTimedOut     
	ErrorCodeHTTPBadStatus   
)

func ErrInvalidConfig(str string) error{

	return nil
}