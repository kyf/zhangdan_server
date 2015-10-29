package main

import (
	"encoding/json"
	"github.com/go-gomail/gomail"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	attach      string = "./attach.txt"
	attachMutex sync.Mutex
)

func response(status bool, msg string) []byte {
	result := map[string]string{
		"status": "ok",
		"msg":    msg,
	}

	if !status {
		result["status"] = "error"
	}

	strresult, _ := json.Marshal(result)
	return strresult
}

func main() {
	http.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		title := r.Form.Get("title")
		body := r.Form.Get("body")

		if strings.EqualFold(title, "") {
			w.Write(response(false, "please pass the title"))
			return
		}

		if strings.EqualFold(body, "") {
			w.Write(response(false, "please pass the body"))
			return
		}

		m := gomail.NewMessage()
		m.SetHeader("From", "kyf456@163.com")
		m.SetHeader("To", "kyf456@163.com", "cora@example.com")
		m.SetHeader("Subject", title)
		m.SetBody("text/plain", body)
		attachMutex.Lock()
		defer attachMutex.Unlock()
		fp, err := os.Create(attach)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
		fp.Write([]byte(body))
		m.Attach(attach)

		d := gomail.NewPlainDialer("smtp.163.com", 25, "kyf456", "1501330364kyf")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
		w.Write(response(true, "data sync success"))
	})
	err := http.ListenAndServe(":2225", nil)
	if err != nil {
		panic(err)
	}

}
