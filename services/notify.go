package services

import (
	"tomelater/global"
	"tomelater/global/log"
	"tomelater/global/time"
)

func Notify() {
	log.Info().Str("service", "Notify").Msg("Start Service")

	go func() {
		ch := global.GetChannel()

		for {
			select {
			case item := <-ch:
				chat.SendTo(item)
			case <-time.After(time.Hour * 24 * 365):
				log.Info().Str("service", "Notify").Msg("Timeout Service")
			}
		}
	}()
}
