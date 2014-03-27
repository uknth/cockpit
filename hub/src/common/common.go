package common

import (
	"strconv"
	"strings"
	"net/smtp"
	"net/mail"
	"encoding/base64"
	"fmt"
	"log"
)

/*
	This function val is for reading the content of _cfile, in case the value isn't present in
	_confMap. It also inserts the value in _confMap, if it isn't present
*/
func val(sec string, opt string) string {
	if sec == "" || opt == "" {
		log.Fatal("Section or Option value sent to config cannot be empty")
	}
	if _confMap[sec+opt] != "" {
		return _confMap[sec+opt]
	} else {
		// cfile is file object from the init.go
		val, err := _cfile.GetRawString(sec, opt)
		if err != nil {
			// Throw exception but let the call proceed
			log.Fatal("Error while getting value for the [" + sec + "/" + opt + "] : (Error) " + err.Error())
			return "err"
		}
		_confMap[sec+opt] = val
		return val
	}
}



func mailer(subject string, content string) bool {
	type Euser struct {
		User   string
		Pass   string
		Server string
		Port   int
	}
	port, _ := strconv.Atoi(val("email", "port"))
	emailUser := &Euser{val("email", "user"), val("email", "password"), val("email", "server"), port}
	auth := smtp.PlainAuth("",
		emailUser.User,
		emailUser.Pass,
		emailUser.Server,
	)
	ml := strings.Split(val("email", "to"), ",")
	for _, val := range ml {
		tos := strings.Split(val, ":")
		from := mail.Address{"Unbxd Monitoring Tool", emailUser.User}
		to := mail.Address{strings.TrimSpace(tos[0]), strings.TrimSpace(tos[1])}
		title := subject
		body := content
		header := make(map[string]string)
		header["From"] = from.String()
		header["To"] = to.String()
		header["Subject"] = func(String string) string {
			// use mail's rfc2047 to encode any string
			addr := mail.Address{String, ""}
			return strings.Trim(addr.String(), " <>")
		}(title)
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		header["Content-Transfer-Encoding"] = "base64"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
		err := smtp.SendMail(
			emailUser.Server+":"+strconv.Itoa(emailUser.Port),
			auth,
			from.Address,
			[]string{to.Address},
			[]byte(message),
		)
		if err != nil {
			log.Print("Unable to send mail" + err.Error())
		}
	}
	return true
}
