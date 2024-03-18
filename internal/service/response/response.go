package response

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(message string) Response {
	return Response{
		Message: message,
	}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
