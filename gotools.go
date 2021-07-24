package gotools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
	err := errors.New("")
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
		err = errC
		val = val2

	}

	return val, err
}
func MakeInt64(deP interface{}) (int64, error) {
	err := errors.New("")
	var val int64
	switch v := deP.(type) {
	case nil:
		val = int64(0.0)
	case int:
		val = int64(v)
	case int32:
		val = int64(v)
	case int64:
		val = int64(v)
	case float32:
		val = int64(v)
	case float64:
		val = int64(v)
	}
	return val, err
}

func makeInt32(deP interface{}) int32 {
	var val int32
	switch v := deP.(type) {
	case nil:
		val = int32(0.0)
	case int:
		val = int32(v)
	case int32:
		val = int32(v)
	case int64:
		val = int32(v)
	case float32:
		val = int32(v)
	case float64:
		val = int32(v)
	}
	return val
}

func makeFloat64(deP interface{}) float64 {
	var val float64
	switch v := deP.(type) {
	case nil:
		val = float64(0.0)
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
			fmt.Println("Erreur conversion string => float", v)

		}
		val = tmp
	}

	return val

}

func sendMail(to string, sub string, msg string) {
	//fmt.Println("Sending mail", to, sub, msg)
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	// Set the sender and recipient.
	c.Mail("admin@a-of")
	c.Rcpt(to)
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	fullMessage := strings.Replace("To: {to}\r\n", "{to}", to, 1)
	fullMessage += strings.Replace("Subject: {sub}\r\n", "{sub}", sub, 1)
	fullMessage += "\r\n" + msg + "\r\n"

	defer wc.Close()
	buf := bytes.NewBufferString(fullMessage)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

func trans(s string, vars ...string) string {
	return s
}

func parseFeed(xml string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(xml)
	return feed

	//func getFeed() *gofeed.Feed {
	//	xml := "https://www.pressherald.com/category/business/feed/"
	//
	//
	//	// fmt.Printf("%+v\n", feed)
	//
	//}
}
func inSlice(n string, l []string) bool {
	for _, s := range l {
		if s == n {
			return true
		}
	}
	return false
}
func cleanString(s string, accent, chars bool) string {
	newS := s
	if accent {
		newS = removeAccents(newS)
	}
	if chars {
		newS = removeSpecialChars(newS)
	}
	return newS

}
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func removeSpecialChars(s string) string {
	repl := strings.NewReplacer("#", "", "@", "", "!", "$", "", "%", "", "^", "", "&", "",
		"*", "", "(", "", ")", " ", "_", ":", "", "[", "", "]", "", "{", "", "}", "", "\\", "", "|", "",
		";", "", "'", "", `"`, "", "?", "", ">", "", "<", "", ",", "", "=", "", "+", "", "`", "", "~", "", ".", "",
	)
	newS := repl.Replace(s)
	return newS

}
func genShittyPasswd() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖabcdefghijklmnopqrstuvwxyzåäöéèÉÈàÀêÊ!@#$%&*=-?()0123456789")
	length := 24
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
func getIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

//

func genClientBsonLink(cli bson.M) string {
	s := ` <a href="/lead-detail/{mid}" target="_blank">{nomCli}</a>`
	mID := cli["_id"].(primitive.ObjectID)
	nom := cli["nom"].(string)
	if _, ok := cli["prenom"].(string); ok {
		nom += " " + cli["prenom"].(string)
	}
	nom = strings.TrimSpace(nom)
	m := mID.Hex()
	np := nom
	np = strings.TrimSpace(np)
	r := strings.NewReplacer("{mid}", m, "{nomCli}", np)
	s = r.Replace(s)
	return s
}
