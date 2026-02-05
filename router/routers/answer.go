package routers

import (

	"strconv"

	"tomelater/global/log"


	"tomelater/controllers/rest"


	"tomelater/models"
	"github.com/gofiber/fiber/v2"
)

// SetupAnswerRoutes sets up routes for answer domain
func SetupAnswerRoutes(group fiber.Router) {

	group.Get("/answer", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.AnswerController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/answer/:id", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.AnswerController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/answer", func(c *fiber.Ctx) error {
		item_ := &models.Answer{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.AnswerController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/answer/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Answer
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.AnswerController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/answer/count", func(c *fiber.Ctx) error {

		var controller rest.AnswerController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/answer", func(c *fiber.Ctx) error {
		item_ := &models.Answer{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.AnswerController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/answer", func(c *fiber.Ctx) error {
		item_ := &models.Answer{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.AnswerController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/answer/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Answer{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.AnswerController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/answer/find/uid/:uid", func(c *fiber.Ctx) error {
		uid_, _ := strconv.ParseInt(c.Params("uid"), 10, 64)
		var controller rest.AnswerController
		controller.Init(c)
		controller.FindByUid(uid_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/answer/find/qid/:qid", func(c *fiber.Ctx) error {
		qid_, _ := strconv.ParseInt(c.Params("qid"), 10, 64)
		var controller rest.AnswerController
		controller.Init(c)
		controller.FindByQid(qid_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}