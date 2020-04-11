package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Krognol/go-wolfram"
	"github.com/christianrondeau/go-wit"
)

var (
	witClient     *wit.Client
	wolframClient *wolfram.Client
)

func main() {

	witClient = wit.NewClient("Y5GBIJSC2VBI5YGQART772JGQAQMRLWQ")
	wolframClient = &wolfram.Client{AppID: "YPAV2U-5X983J5VXJ"}
	reader := bufio.NewReader(os.Stdin)

	for true {
		msg, _ := reader.ReadString('\n')
		handleMessage(msg)

	}

}

func handleMessage(msg string) {

	var (
		confidenceThreshold = 0.5
		topEntity           wit.MessageEntity
		topEntityKey        string
	)

	resp, err := witClient.Message(msg)
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
	replyToUser(topEntityKey, topEntity)
}

func replyToUser(topEntityKey string, entity wit.MessageEntity) {

	var myText string
	switch strings.ToLower(topEntityKey) {
	case "greetings":
		myText = "Hi, How can i help you?"
		fmt.Println(myText)
		return
	case "bye":
		myText = "Thank you, Have a nice day."
		fmt.Println(myText)
		return
	case "thanks":
		myText = "Its my pleasure. Please let me know if i can help yu with anything else."
		fmt.Println(myText)
		return
	case "wolfram_search_query":
		res, err := wolframClient.GetShortAnswerQuery(entity.Value.(string), wolfram.Metric, 1000)
		if err != nil {
			log.Printf("unable to get wolfram result: %v", err)
			return
		}
		fmt.Println(res)
		return
	default:
		myText = "Hi it is out of my scope. sorry"
		fmt.Println(myText)
		return
	}

}
