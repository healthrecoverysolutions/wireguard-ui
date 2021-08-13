package emailer

import (
	"bytes"
	"fmt"
	"net/textproto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/jordan-wright/email"
)

type SESApiMail struct {
	region   string
	from     string
}

func NewSESApiMail(region, from string) *SESApiMail {
	ans := SESApiMail{region: region, from: from}
	return &ans
}

func (o *SESApiMail) Send(toName string, to string, subject string, content string, attachments []Attachment) error {

	sess, err := session.NewSession(&aws.Config{
		Region:aws.String(o.region)},
	)
	svc := ses.New(sess)

	e := &email.Email {
		To: []string{to},
		From: o.from,
		Subject: subject,
		HTML: []byte(content),
		Headers: textproto.MIMEHeader{},
	}
	for _, attachment := range attachments {
		attachmentReader := bytes.NewReader(attachment.Data)
		_, _ = e.Attach(attachmentReader, attachment.Name, attachment.ContentType)
	}

	rawEmail, _ := e.Bytes()
	message := ses.RawMessage{ Data: rawEmail}

	source := aws.String(o.from)
	destinations := []*string{aws.String(to)}
	input := ses.SendRawEmailInput{Source: source, Destinations: destinations, RawMessage: &message}
	result, err := svc.SendRawEmail(&input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Email Sent to address: " + to)
	fmt.Println(result.String())

	return nil
}
