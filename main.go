// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var boyMapping map[string]string
var girlMapping map[string]string

func main() {
	var err error
	boyMapping = make(map[string]string)
	girlMapping = make(map[string]string)
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}


func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {

		if event.Type == linebot.EventTypeMessage {

			switch message := event.Message.(type) {
			
			case *linebot.TextMessage:

				if strings.Contains(message.Text, "#") {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.UserID + " : " + message.Text+" 裡面有#")).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "隨機聊天" {
					leftBtn := linebot.NewMessageTemplateAction("我是男生", "-我是男生，開始尋找配對中...")
					rightBtn := linebot.NewMessageTemplateAction("我是女生", "-我是女生，開始尋找配對中...")
					template := linebot.NewConfirmTemplate("隨機聊天", leftBtn, rightBtn)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("開始隨機聊天。", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "-我是男生，開始尋找配對中..." {
					
					boyMapping[event.Source.UserID] := ""

					for girlId := range girlMapping {
					    
					    if girlMapping[girlId] == "" {
					    	girlMapping[girlId] = event.Source.UserID
					    	boyMapping[event.Source.UserID] = girlId
					    }
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("開始隨機聊天。", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "-我是女生，開始尋找配對中..." {
					

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("開始隨機聊天。", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if strings.Contains(message.Text, "**") {

					if _, err := bot.PushMessage("U19cf2d03ac94f9ecc9c3ec3664b80f6d", linebot.NewTextMessage("私訊你" + event.Source.UserID)).Do(); err != nil {
					log.Print(err)
					}
					if _, err := bot.PushMessage("U19cf2d03ac94f9ecc9c3ec3664b80f6d", linebot.NewTextMessage("私訊你" + event.Source.UserID)).Do(); err != nil {
					log.Print(err)
					}
				}
				
			}
		}
	}
}
