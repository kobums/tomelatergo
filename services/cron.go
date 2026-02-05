package services

import (
	"tomelater/global"
	"tomelater/global/log"
	"tomelater/global/time"
	"tomelater/models"

	"github.com/robfig/cron"
)

var _cron *cron.Cron

func Cron() {
	log.Info().Str("service", "Cron").Msg("Start Service")

	go func() {
		InitCron()
		ch := global.GetCronChannel()

		for {
			select {
			case <-ch:
				RestartCron()
			case <-time.After(time.Hour * 24 * 365):
				log.Info().Str("service", "Cron").Msg("Timeout Service")
			}
		}
	}()
}

func InitCron() {
	_cron = cron.New()

	// instance := setting.GetInstance()
	// team := instance.Setting("cron.crawling.team")
	// area := instance.Setting("cron.crawling.area")
	// grade := instance.Setting("cron.crawling.grade")
	// goldenage := instance.Setting("cron.crawling.goldenage")
	// represent := instance.Setting("cron.crawling.represent")
	// game := instance.Setting("cron.crawling.game")
	// match := instance.Setting("cron.crawling.match")

	err := _cron.AddFunc("0 0 0 * * *", log.Rotate)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	err = _cron.AddFunc("0 0,30 * * * *", boardCounter)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// err = _cron.AddFunc(team, crawler.SyncTeam)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(area, crawler.SyncArea)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(grade, crawler.SyncGrade)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(goldenage, crawler.SyncGoldenage)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(represent, crawler.SyncRepresent)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(game, crawler.SyncGame)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// err = _cron.AddFunc(match, crawler.SyncMatch)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	_cron.Start()
}

func RestartCron() {
	log.Info().Str("service", "Cron").Msg("Restart Service")

	_cron.Stop()

	InitCron()
}

func boardCounter() {
	conn := models.NewConnection()
	defer conn.Close()

	// boardManager := models.NewBoardManager(conn)
	// freeshareManager := models.NewFreeshareManager(conn)
	// teampromotionManager := models.NewTeampromotionManager(conn)
	// recommendhotelManager := models.NewRecommendhotelManager(conn)
	// certicompanyManager := models.NewCerticompanyManager(conn)
	// bignewsManager := models.NewBignewsManager(conn)

	// {
	// 	items := boardManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		boardManager.IncreaseCountById(1, v.Id)
	// 	}
	// }

	// {
	// 	items := freeshareManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		freeshareManager.IncreaseCountById(1, v.Id)
	// 	}
	// }

	// {
	// 	items := teampromotionManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		teampromotionManager.IncreaseCountById(1, v.Id)
	// 	}
	// }
	// {
	// 	items := recommendhotelManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		recommendhotelManager.IncreaseCountById(1, v.Id)
	// 	}
	// }
	// {
	// 	items := certicompanyManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		certicompanyManager.IncreaseCountById(1, v.Id)
	// 	}
	// }
	// {
	// 	items := bignewsManager.Find([]interface{}{
	// 		models.Where{Column: "count", Value: 200, Compare: "<="},
	// 		models.Ordering("rand()"),
	// 		models.Limit(1),
	// 	})

	// 	for _, v := range items {
	// 		bignewsManager.IncreaseCountById(1, v.Id)
	// 	}
	// }
}
