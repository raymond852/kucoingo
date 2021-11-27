package kucoingo

import "fmt"

type ResponseError struct {
	Body       []byte
	StatusCode int
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("statusCode=%d body=%s", e.StatusCode, string(e.Body))
}

type BadRequestError struct {
	MissingFields   []string
	IncorrectFields []string
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("missing fields = %+v, incorrect fields = %+v", e.MissingFields, e.IncorrectFields)
}
