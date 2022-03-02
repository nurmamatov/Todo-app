package repo

// Email ...
type Email struct {
	ID             string
	Subject        string
	Body           string
	RecipientEmail string
	Phone          string
}
type Sms struct {
	Text string
	Num  string
}

// SendStorageI ...
type SendStorageI interface {
	MakeSent(ID string) error
	Send(subject, body string, status bool, recipients string) error
	SendSms(phone_number string, text string) error
}
