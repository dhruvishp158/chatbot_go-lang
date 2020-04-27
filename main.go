package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Krognol/go-wolfram"
	"github.com/christianrondeau/go-wit"
)

var tpl *template.Template

var (
	witClient     *wit.Client
	wolframClient *wolfram.Client
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
func main() {

	witClient = wit.NewClient("Y5GBIJSC2VBI5YGQART772JGQAQMRLWQ")
	wolframClient = &wolfram.Client{AppID: "YPAV2U-5X983J5VXJ"}

	http.HandleFunc("/", index)
	http.HandleFunc("/process", handleMessage)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	chatInput := r.FormValue("chat")

	var (
		confidenceThreshold = 0.5
		topEntity           wit.MessageEntity
		topEntityKey        string
	)

	resp, err := witClient.Message(chatInput)
	if err != nil {
		log.Printf("unable to get wit.ai responses: %v", err)
		return
	}

	for entityKey, entityList := range resp.Entities {
		for _, entity := range entityList {
			if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
				topEntity = entity
				topEntityKey = entityKey

			}
		}
	}
	out := replyToUser(topEntityKey, topEntity, w, r)

	d := struct {
		Input, Output string
	}{
		Input:  chatInput,
		Output: out,
	}
	tpl.ExecuteTemplate(w, "index.gohtml", d)

}

func replyToUser(topEntityKey string, entity wit.MessageEntity, w http.ResponseWriter, r *http.Request) string {

	var res string
	switch strings.ToLower(topEntityKey) {
	case "greetings":
		res = "Hi, How can i help you?"

	case "bye":
		res = "Thank you, Have a nice day."

	case "thanks":
		res = "Its my pleasure. Please let me know if i can help yu with anything else."

	case "wolfram_search_query":
		r, err := wolframClient.GetShortAnswerQuery(entity.Value.(string), wolfram.Metric, 1000)
		if err != nil {
			log.Printf("unable to get wolfram result: %v", err)
		}
		res = r

	default:
		res = "Hi it is out of my scope. sorry"

	}

	return res

}
