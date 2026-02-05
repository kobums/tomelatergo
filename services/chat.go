package services

import (
	"fmt"
	"time"
	"tomelater/global"
	"tomelater/global/log"
	"tomelater/global/setting"

	"github.com/YiYuhki/ikisocket"
)

type ChatSession struct {
	Id    int64
	UUID  string
	Valid bool
	Date  time.Time
}

type ChatService struct {
	Use     bool
	Clients map[string]ChatSession
}

type MessageObject struct {
	Data string `json:"data"`
	From string `json:"from"`
	To   string `json:"to"`
}

var chat ChatService

func (p *ChatService) SendTo(notify global.Notify) {
	id := notify.Id
	message := notify.Message
	uuid := notify.UUID

	for _, v := range p.Clients {
		if id == v.Id {
			if message != global.SessionTimeout {
				uuid = v.UUID
			}

			err := ikisocket.EmitTo(uuid, []byte(global.GetMessage(message)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func Chat() {
	log.Info().Str("service", "Session").Msg("Start Service")

	chat.Use = true
	chat.Clients = make(map[string]ChatSession)

	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		err := ep.Kws.EmitTo(ep.Kws.UUID, []byte("pong"))
		if err != nil {
			fmt.Println(err)
		}
	})

	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		// id := ep.Kws.GetAttribute("id").(int64)

		value, exist := chat.Clients[ep.Kws.UUID]
		if !exist {
			err := ep.Kws.EmitTo(ep.Kws.UUID, []byte(global.GetMessage(global.SessionTimeout)))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		instance := setting.GetInstance()
		timeout := instance.SettingInt("user.login.timeout")

		if timeout > 0 {
			last := value.Date.Add(time.Second * time.Duration(timeout))
			if last.Before(time.Now()) {
				err := ep.Kws.EmitTo(ep.Kws.UUID, []byte(global.GetMessage(global.SessionTimeout)))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}

		err := ep.Kws.EmitTo(ep.Kws.UUID, []byte("pong"))
		if err != nil {
			fmt.Println(err)
		}

		if string(ep.Data) == "ping" {
			flag := false

			if !exist {
				flag = true
			} else {
				if instance.Setting("user.login.single") == "Y" {
					if exist {
						if !value.Valid {
							flag = true
						}
					} else {
						flag = true
					}

				}
			}

			if flag {
				err := ep.Kws.EmitTo(ep.Kws.UUID, []byte(global.GetMessage(global.SessionTimeout)))
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		} else {
			value.Date = time.Now()
			chat.Clients[value.UUID] = value
		}
	})

	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		// id := ep.Kws.GetAttribute("id").(int64)

		delete(chat.Clients, ep.Kws.UUID)
	})

	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		// id := ep.Kws.GetAttribute("id").(int64)

		delete(chat.Clients, ep.Kws.UUID)
	})

	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		//id := ep.Kws.GetAttribute("id").(int64)
		delete(chat.Clients, ep.Kws.UUID)
	})
}
