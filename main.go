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

				if message.Text == "選單" {
					leftBtn := linebot.NewMessageTemplateAction("停止聊天", "*已停止隨機聊天功能。")
					rightBtn := linebot.NewMessageTemplateAction("下一位", "*尋找下一位聊天對象中...")
					template := linebot.NewConfirmTemplate("聊天設定", leftBtn, rightBtn)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("聊天設定", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "開始聊天" {
					leftBtn := linebot.NewMessageTemplateAction("我是男生", "*我是男生，開始尋找配對中...")
					rightBtn := linebot.NewMessageTemplateAction("我是女生", "*我是女生，開始尋找配對中...")
					template := linebot.NewConfirmTemplate("請選擇", leftBtn, rightBtn)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("開始隨機聊天。", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "狀態" {
					leftBtn := linebot.NewMessageTemplateAction("我是男生", "*我是男生，我要查詢我目前在聊天室的狀態...")
					rightBtn := linebot.NewMessageTemplateAction("我是女生", "*我是女生，我要查詢我目前在聊天室的狀態...")
					template := linebot.NewConfirmTemplate("請選擇", leftBtn, rightBtn)

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("狀態查詢", template)).Do(); err != nil {
					log.Print(err)
					}
				} else if message.Text == "*我是男生，開始尋找配對中..." {
					
					boyMapping[event.Source.UserID] = "wait"

					for girlId := range girlMapping {
					    
					    if girlMapping[girlId] == "wait" {
					    	girlMapping[girlId] = event.Source.UserID
					    	boyMapping[event.Source.UserID] = girlId
					    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
							log.Print(err)
							}
							if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
							log.Print(err)
							}
							break
					    }
					}

				} else if message.Text == "*我是女生，開始尋找配對中..." {

					girlMapping[event.Source.UserID] = "wait"

					for boyId := range boyMapping {
					    
					    if boyMapping[boyId] == "wait" {
					    	boyMapping[boyId] = event.Source.UserID
					    	girlMapping[event.Source.UserID] = boyId
					    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
							log.Print(err)
							}
							if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
							log.Print(err)
							}
							break
					    }
					}
					
				} else if message.Text == "*我是男生，我要查詢我目前在聊天室的狀態..." {

					if boyMapping[event.Source.UserID] == "wait" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("尋找配對中")).Do(); err != nil {
							log.Print(err)
						}
					} else if boyMapping[event.Source.UserID] == "" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("目前已停止聊天功能")).Do(); err != nil {
							log.Print(err)
						}
					} else if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已配對聊天(只是對方不理你)")).Do(); err != nil {
							log.Print(err)
						}
					}
					
				} else if message.Text == "*我是女生，我要查詢我目前在聊天室的狀態..." {

					if girlMapping[event.Source.UserID] == "wait" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("尋找配對中")).Do(); err != nil {
							log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] == "" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("目前已停止聊天功能")).Do(); err != nil {
							log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已配對聊天(只是對方不理你)")).Do(); err != nil {
							log.Print(err)
						}
					}
					
				} else if message.Text == "*已停止隨機聊天功能。" {

					if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(girlMapping[event.Source.UserID], linebot.NewTextMessage("對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						boyMapping[girlMapping[event.Source.UserID]] = "wait"
						var boyId = girlMapping[event.Source.UserID]
						girlMapping[event.Source.UserID] = ""

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" {
						    	girlMapping[girlId] = boyId
						    	boyMapping[boyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

					}

					if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(boyMapping[event.Source.UserID], linebot.NewTextMessage("對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						girlMapping[boyMapping[event.Source.UserID]] = "wait"
						var girlId = boyMapping[event.Source.UserID]
						boyMapping[event.Source.UserID] = ""

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" {
						    	boyMapping[boyId] = girlId
						    	girlMapping[girlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}
					
				} else if message.Text == "*尋找下一位聊天對象中..." {

					if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
					
						if _, err := bot.PushMessage(girlMapping[event.Source.UserID], linebot.NewTextMessage("對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var tempGirlId = event.Source.UserID
						var tempBoyId = girlMapping[event.Source.UserID]
						boyMapping[girlMapping[event.Source.UserID]] = "wait"
						girlMapping[event.Source.UserID] = "wait"

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" {
						    	boyMapping[boyId] = tempGirlId
						    	girlMapping[tempGirlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempGirlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" {
						    	girlMapping[girlId] = tempBoyId
						    	boyMapping[tempBoyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempBoyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}

					if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(boyMapping[event.Source.UserID], linebot.NewTextMessage("對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var tempGirlId = boyMapping[event.Source.UserID]
						var tempBoyId = event.Source.UserID
						girlMapping[boyMapping[event.Source.UserID]] = "wait"
						boyMapping[event.Source.UserID] = "wait"

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" {
						    	boyMapping[boyId] = tempGirlId
						    	girlMapping[tempGirlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempGirlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" {
						    	girlMapping[girlId] = tempBoyId
						    	boyMapping[tempBoyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempBoyId, linebot.NewTextMessage("配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}
					
				} else if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {

					if _, err := bot.PushMessage(boyMapping[event.Source.UserID], linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
					}
				} else if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {

					if _, err := bot.PushMessage(girlMapping[event.Source.UserID], linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
					}
				} else if girlMapping[event.Source.UserID] == "wait" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("尚未配對成功，請稍候。")).Do(); err != nil {
						log.Print(err)
					}
				} else if boyMapping[event.Source.UserID] == "wait" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("尚未配對成功，請稍候。")).Do(); err != nil {
						log.Print(err)
					}
				}
				
			}
		}
	}
}
