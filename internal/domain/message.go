package domain


type Request struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
}

type Response struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
}