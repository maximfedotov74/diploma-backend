package fall

type Error interface {
	Status() int
	Message() string
}

type AppErr struct {
	MessageText string `json:"message" example:"some error" validate:"required"`
	StatusCode  int    `json:"status" example:"404" validate:"required"`
}

func NewErr(m string, s int) Error {
	return AppErr{MessageText: m, StatusCode: s}
}

func (err AppErr) Status() int {
	return err.StatusCode
}

func (err AppErr) Message() string {
	return err.MessageText
}

func ServerError(m string) AppErr {
	return AppErr{
		MessageText: "Internal Server Error \n" + m,
		StatusCode:  STATUS_INTERNAL_ERROR,
	}
}
