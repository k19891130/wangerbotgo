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
	"strconv"
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
				} else if message.Text == "#31#Profit" {

					for girlId := range girlMapping {
						res, err := bot.GetProfile(girlId).Do();
						if err != nil {
						    if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("性別 : 女\n" + "名字 : " + res.DisplayName + "\n照片 : " + res.PictureUrl + "\nID : " + res.UserID + "\n個人狀態 : " + res.StatusMessage)).Do(); err != nil {
							log.Print(err)
							}
						}
					}

					for boyId := range boyMapping {
						res, err := bot.GetProfile(boyId).Do();
						if err != nil {
						    if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("性別 : 男\n" + "名字 : " + res.DisplayName + "\n照片 : " + res.PictureUrl + "\nID : " + res.UserID + "\n個人狀態 : " + res.StatusMessage)).Do(); err != nil {
							log.Print(err)
							}
						}
					}

				} else if message.Text == "開始聊天" {

					if boyMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*你已在男生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*你已在女生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else {
						leftBtn := linebot.NewMessageTemplateAction("我是男生", "*我是男生，開始尋找配對中...")
						rightBtn := linebot.NewMessageTemplateAction("我是女生", "*我是女生，開始尋找配對中...")
						template := linebot.NewConfirmTemplate("請選擇", leftBtn, rightBtn)

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("*開始隨機聊天。", template)).Do(); err != nil {
						log.Print(err)
						}
					}
				} else if message.Text == "狀態" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*目前聊天室中有 : \n" + strconv.Itoa(len(boyMapping)) +" 個男生\n" + strconv.Itoa(len(girlMapping)) +" 個女生")).Do(); err != nil {
						log.Print(err)
					}

					if boyMapping[event.Source.UserID] == "wait" {
						if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*你的狀態為 : \n尋找配對中")).Do(); err != nil {
							log.Print(err)
						}
					} else if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*你的狀態為 : \n已配對聊天(只是對方不理你)")).Do(); err != nil {
							log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] == "wait" {
						if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*你的狀態為 : \n尋找配對中")).Do(); err != nil {
							log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
						if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*你的狀態為 : \n已配對聊天(只是對方不理你)")).Do(); err != nil {
							log.Print(err)
						}
					} else {
						if _, err = bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*你的狀態為 : \n尚未開啟隨機聊天功能，請輸入\"開始聊天\"來加入聊天序列。")).Do(); err != nil {
							log.Print(err)
						}
					}
				} else if message.Text == "*我是男生，開始尋找配對中..." {

					if boyMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*你已在男生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*妳已在女生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else {
						boyMapping[event.Source.UserID] = "wait"

						for girlId := range girlMapping {
						    
						    if girlMapping[girlId] == "wait" {
						    	girlMapping[girlId] = event.Source.UserID
						    	boyMapping[event.Source.UserID] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}

				} else if message.Text == "*我是女生，開始尋找配對中..." {

					if boyMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*你已在男生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else if girlMapping[event.Source.UserID] != "" {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*妳已在女生的等待序列中。")).Do(); err != nil {
						log.Print(err)
						}
					} else {
						girlMapping[event.Source.UserID] = "wait"

						for boyId := range boyMapping {
						    
						    if boyMapping[boyId] == "wait" {
						    	boyMapping[boyId] = event.Source.UserID
						    	girlMapping[event.Source.UserID] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(event.Source.UserID, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}					
				} else if message.Text == "*已停止隨機聊天功能。" {

					if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(girlMapping[event.Source.UserID], linebot.NewTextMessage("*對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var boyId = girlMapping[event.Source.UserID]
						boyMapping[boyId] = "wait"
						girlMapping[event.Source.UserID] = ""

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" {
						    	girlMapping[girlId] = boyId
						    	boyMapping[boyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

					}

					if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(boyMapping[event.Source.UserID], linebot.NewTextMessage("*對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var girlId = boyMapping[event.Source.UserID]
						girlMapping[girlId] = "wait"
						boyMapping[event.Source.UserID] = ""

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" {
						    	boyMapping[boyId] = girlId
						    	girlMapping[girlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					}
					delete(girlMapping, event.Source.UserID)
					delete(boyMapping, event.Source.UserID)
					
				} else if message.Text == "*尋找下一位聊天對象中..." {

					if girlMapping[event.Source.UserID] != "" && len(girlMapping[event.Source.UserID]) > 10 {
					
						if _, err := bot.PushMessage(girlMapping[event.Source.UserID], linebot.NewTextMessage("*對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var tempGirlId = event.Source.UserID
						var tempBoyId = girlMapping[event.Source.UserID]
						boyMapping[girlMapping[event.Source.UserID]] = "wait"
						girlMapping[event.Source.UserID] = "wait"

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" && tempBoyId != boyId {
						    	boyMapping[boyId] = tempGirlId
						    	girlMapping[tempGirlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempGirlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" && girlId != tempGirlId {
						    	girlMapping[girlId] = tempBoyId
						    	boyMapping[tempBoyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempBoyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					} else if boyMapping[event.Source.UserID] != "" && len(boyMapping[event.Source.UserID]) > 10 {
						
						if _, err := bot.PushMessage(boyMapping[event.Source.UserID], linebot.NewTextMessage("*對方已離開，重新尋找對象中。")).Do(); err != nil {
							log.Print(err)
						}
						var tempGirlId = boyMapping[event.Source.UserID]
						var tempBoyId = event.Source.UserID
						girlMapping[boyMapping[event.Source.UserID]] = "wait"
						boyMapping[event.Source.UserID] = "wait"

						for boyId := range boyMapping {
					    
						    if boyMapping[boyId] == "wait" && tempBoyId != boyId {
						    	boyMapping[boyId] = tempGirlId
						    	girlMapping[tempGirlId] = boyId
						    	if _, err := bot.PushMessage(boyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempGirlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}

						for girlId := range girlMapping {
					    
						    if girlMapping[girlId] == "wait" && girlId != tempGirlId {
						    	girlMapping[girlId] = tempBoyId
						    	boyMapping[tempBoyId] = girlId
						    	if _, err := bot.PushMessage(girlId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								if _, err := bot.PushMessage(tempBoyId, linebot.NewTextMessage("*配對成功，請開始聊天。")).Do(); err != nil {
								log.Print(err)
								}
								break
						    }
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*您已在等待序列中。")).Do(); err != nil {
							log.Print(err)
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

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*尚未配對成功，請稍候。")).Do(); err != nil {
						log.Print(err)
					}
				} else if boyMapping[event.Source.UserID] == "wait" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*尚未配對成功，請稍候。")).Do(); err != nil {
						log.Print(err)
					}
				} else if girlMapping[event.Source.UserID] == "" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*尚未開啟聊天功能，請輸入\"開始聊天\"。")).Do(); err != nil {
						log.Print(err)
					}
				} else if boyMapping[event.Source.UserID] == "" {

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*尚未開啟聊天功能，請輸入\"開始聊天\"。")).Do(); err != nil {
						log.Print(err)
					}
				}
				
			}
		} else {
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("*抱歉，目前只開放純文字聊天。")).Do(); err != nil {
				log.Print(err)
			}
		}
	}
}

