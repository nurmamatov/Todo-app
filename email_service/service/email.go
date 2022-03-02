package service

import (
	"app/config"
	pb "app/genproto/email"
	"app/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gomail "gopkg.in/gomail.v2"
)

// SendService ...
type SendService struct {
	storage storage.I
	conf    config.Config
}
type SMSRequestBody struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// NewSendService ...
func NewSendService(db *sqlx.DB, cfg config.Config) *SendService {
	return &SendService{storage: storage.NewStoragePg(db), conf: cfg}
}

//Send ...
func (s *SendService) SendEmail(ctx context.Context, req *pb.Email) (*pb.Empty, error) {
	
	statuss := true
	err := s.SendToEmail(req.Subject, req.Body, req.Email)
	log.Print(err)
	if err != nil {
		statuss = false
		err = s.storage.SendEmail().Send(req.Subject, req.Body, statuss, req.Email)
		if err != nil {
			return &pb.Empty{}, status.Error(codes.Internal, "Internal server error")
		}

	} else {
		statuss = true
		err = s.storage.SendEmail().Send(req.Subject, req.Body, statuss, req.Email)
		if err != nil {
			return nil, err
		}
	}
	return &pb.Empty{}, nil
}

func (s *SendService) SendSms(ctx context.Context, req *pb.Sms) (*pb.Status, error) {
	err := s.SendToSms(req.Phone, req.Body)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Status: "Ok!"}, nil
}


func (s *SendService) SendToEmail(subject string, body string, email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.conf.EmailFromHeader)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email to
	d := gomail.NewPlainDialer(s.conf.SMTPHost, s.conf.SMTPPort, s.conf.SMTPUser, s.conf.SMTPUserPass)

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
		panic(err)
	}
	log.Print("Sent")
	return nil
}

func (s *SendService) SendToSms(phone_number string, text string) error {
	body := SMSRequestBody{}
	body.APIKey = "ae54c246"
	body.APISecret = "6oUDWGyRAEFPcegb"
	body.To = phone_number
	body.From = "Khusniddin"
	body.Text = text

	smsBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("https://rest.nexmo.com/sms/json", "application/json", bytes.NewBuffer(smsBody))
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = s.storage.SendEmail().SendSms(phone_number, text)
	if err != nil {
		return err
	}
	fmt.Println(string(respBody))
	resp.Body.Close()
	return nil
}
