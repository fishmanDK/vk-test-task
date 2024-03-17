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

//func ValidationError(errs validator.ValidationErrors) Response {
//	var errMsgs []string
//
//	for _, err := range errs {
//		switch err.ActualTag() {
//		case "required":
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
//		case "url":
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
//		default:
//			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
//		}
//	}
//
//	return Response{
//		Status: StatusError,
//		Error:  strings.Join(errMsgs, ", "),
//	}
//}
