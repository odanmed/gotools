package gotools

import (
	"bytes"
	"encoding/json"

	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ResultS struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ReturnRes(res ResultS, code int) {
	jRes, _ := json.Marshal(res)
	fmt.Println(string(jRes))
	os.Exit(code)

}

func MakeInt(deP interface{}) (int, error) {
	var val int

	switch v := deP.(type) {
	case nil:
		val = 0
	case int:
		val = v
	case int32:
		val = int(v)
	case int64:
		val = int(v)
	case float32:
		val = int(v)
	case float64:
		val = int(v)
	case string:
		val2, errC := strconv.Atoi(v)
		if errC != nil {
			return 0, errC
		}
		val = val2

	}

	return val, nil
}
func MakeInt64(deP interface{}) (int64, error) {

	var val int64
	switch v := deP.(type) {
	case nil:
		val = int64(0.0)
	case int:
		val = int64(v)
	case int32:
		val = int64(v)
	case int64:
		val = v
	case float32:
		val = int64(v)
	case float64:
		val = int64(v)
	}
	return val, nil
}

func MakeInt32(deP interface{}) (int32, error) {
	var val int32
	switch v := deP.(type) {
	case nil:
		val = int32(0.0)
	case int:
		val = int32(v)
	case int32:
		val = v
	case int64:
		val = int32(v)
	case float32:
		val = int32(v)
	case float64:
		val = int32(v)
	}
	return val, nil
}

func MakeFloat64(deP interface{}) (float64, error) {
	var val float64
	switch v := deP.(type) {
	case nil:
		val = 0.0
	case int:
		val = float64(v)
	case int32:
		val = float64(v)
	case int64:
		val = float64(v)
	case float32:
		val = float64(v)
	case float64:
		val = v
	case string:
		tmp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		val = tmp
	}

	return val, nil

}
func SendHtmlMailLocalhost(to string, from string, sub string, msg string) error {
	//fmt.Println("Sending mail", to, sub, msg)

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		return err
	}
	defer c.Close()
	// Set the sender and recipient.
	err = c.Mail(from)
	if err != nil {
		return err
	}
	err = c.Rcpt(to)
	if err != nil {
		return err
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	fullMessage := strings.Replace("To: {to}\r\n", "{to}", to, 1)
	fullMessage += strings.Replace("Subject: {sub}\r\n", "{sub}", sub, 1)
	fullMessage += `Mime-Version: 1.0;` + "\r\n"
	fullMessage += `Content-Type: text/html; charset="utf-8";` + "\r\n"
	fullMessage += `Content-Transfer-Encoding: 7bit;`
	fullMessage += "\r\n<html><body>" + msg + "</body></html>\r\n"

	defer wc.Close()
	buf := bytes.NewBufferString(fullMessage)
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}
	return nil
}
func SendTextMailLocalhost(to string, from string, sub string, msg string) error {
	//fmt.Println("Sending mail", to, sub, msg)

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		return err
	}
	defer c.Close()
	// Set the sender and recipient.
	err = c.Mail(from)
	if err != nil {
		return err
	}
	err = c.Rcpt(to)
	if err != nil {
		return err
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	fullMessage := strings.Replace("To: {to}\r\n", "{to}", to, 1)
	fullMessage += strings.Replace("Subject: {sub}\r\n", "{sub}", sub, 1)

	fullMessage += "\r\n" + msg + "\r\n"

	defer wc.Close()
	buf := bytes.NewBufferString(fullMessage)
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}
	return nil
}
func InSlice(n string, l []string) bool {
	for _, s := range l {
		if s == n {
			return true
		}
	}
	return false
}
func InIntSlice(n int, l []int) bool {
	for _, s := range l {
		if s == n {
			return true
		}
	}
	return false
}
func CleanString(s string, accent, chars bool) string {
	newS := s
	if accent {
		newS = RemoveAccents(newS)
	}
	if chars {
		newS = RemoveSpecialChars(newS)
	}
	return newS

}
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func RemoveSpecialChars(s string) string {
	repl := strings.NewReplacer("#", "", "@", "", "!", "$", "", "%", "", "^", "", "&", "",
		"*", "", "(", "", ")", " ", "_", ":", "", "[", "", "]", "", "{", "", "}", "", "\\", "", "|", "",
		";", "", "'", "", `"`, "", "?", "", ">", "", "<", "", ",", "", "=", "", "+", "", "`", "", "~", "", ".", "",
	)
	newS := repl.Replace(s)
	return newS

}
func GenPasswordALaCon() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖabcdefghijklmnopqrstuvwxyzåäöéèÉÈàÀêÊ!@#$%&*=-?()0123456789")
	length := 24
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
func GetIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
