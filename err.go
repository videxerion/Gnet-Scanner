package main

type EmptyResponse struct {
	message string
}

func (e *EmptyResponse) Error() string {
	return e.message
}
