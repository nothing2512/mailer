# Mailer
Go-Mailer Plugins

## Installation
```bash
go get -u github.com/nothing2512/mailer
```

## Example Usage
```go

package main

import (
    "github.com/nothing2512/mailer"
)

type Data struct {
    Name string `json:"name"`
}

func main() {
    m, err := mailer.Init("email", "pass", "host", "port")
	if err != nil {
		panic(err)
	}
	data := Data{"Fulan"}

	m.Subject("subject")

	m.Recipients("mail@gmail.com")
	m.Cc("mail@gmail.com", "mail@gmail.com")
	m.Bcc("mail@gmail.com", "mail@gmail.com")

	m.SetHTMLFile("template.html", data)
	m.AttachFile("certificate.pdf", []byte{})

	err = m.Send()
	if err != nil {
		panic(err)
	}
}
```
