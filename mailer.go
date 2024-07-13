package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"strings"
)

type Mailer struct {
	email    string
	password string
	host     string
	port     string
	client   *smtp.Client

	buffer *bytes.Buffer
	writer *multipart.Writer

	recipients []string
	subject    string

	cc  []string
	bcc []string
}

func Init(email, password, host, port string) (*Mailer, error) {
	auth := smtp.PlainAuth("", email, password, host)
	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return nil, err
	}
	if err := client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return nil, err
	}
	if err := client.Auth(auth); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	return &Mailer{
		email:    email,
		password: password,
		host:     host,
		port:     port,
		client:   client,
		writer:   writer,
		buffer:   &buffer,
	}, nil
}

func (m *Mailer) Cc(emails ...string) {
	m.cc = emails
}

func (m *Mailer) Bcc(emails ...string) {
	m.bcc = emails
}

func (m *Mailer) getHeaders() string {
	headers := make(map[string]string)
	headers["From"] = m.email
	headers["To"] = strings.Join(m.recipients, ", ")
	if len(m.cc) > 0 {
		headers["Cc"] = strings.Join(m.cc, ", ")
	}
	headers["Subject"] = m.subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + m.writer.Boundary()

	h := ""

	for k, v := range headers {
		h += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return h + ";\r\n\n"
}

func (m *Mailer) Subject(subject string) {
	m.subject = subject
}

func (m *Mailer) Recipients(emails ...string) {
	m.recipients = emails
}

func (m *Mailer) SetText(text string) error {
	part, err := m.writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	part.Write([]byte(text))
	return nil
}

func (m *Mailer) SetHTML(content string) error {
	textPart, err := m.writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	textPart.Write([]byte(content))
	return nil
}

func (m *Mailer) SetHTMLFile(file string, data interface{}) error {

	var body bytes.Buffer
	t, err := template.ParseFiles(file)
	if err != nil {
		panic(err)
	}
	_ = t.Execute(&body, data)

	textPart, err := m.writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=\"utf-8\""},
	})
	if err != nil {
		return err
	}
	textPart.Write(body.Bytes())
	return nil
}

func (m *Mailer) AttachFile(filename string, data []byte) error {
	part, err := m.writer.CreatePart(map[string][]string{
		"Content-Disposition":       {"attachment; filename=\"" + filename + "\""},
		"Content-Type":              {"application/octet-stream"},
		"Content-Transfer-Encoding": {"base64"},
	})
	if err != nil {
		return err
	}

	b64 := base64.NewEncoder(base64.StdEncoding, part)
	if _, err = b64.Write(data); err != nil {
		panic(err)
	}
	b64.Close()
	return nil
}

func (m *Mailer) Send() error {
	m.writer.Close()
	headers := m.getHeaders()
	content := m.buffer.String()

	m.client.Mail(m.email)

	for _, x := range append(m.recipients, append(m.cc, m.bcc...)...) {
		m.client.Rcpt(x)
	}

	w, err := m.client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(headers + content))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *Mailer) Close() {
	m.client.Close()
}
