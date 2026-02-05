package routers

import (

	"strconv"

	"tomelater/global/log"


	"tomelater/controllers/rest"


	"tomelater/models"
	"github.com/gofiber/fiber/v2"
)

// SetupQuestionRoutes sets up routes for question domain
func SetupQuestionRoutes(group fiber.Router) {

	group.Get("/question", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.QuestionController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/question/:id", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.QuestionController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/question", func(c *fiber.Ctx) error {
		item_ := &models.Question{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.QuestionController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/question/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.Question
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.QuestionController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/question/count", func(c *fiber.Ctx) error {

		var controller rest.QuestionController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/question", func(c *fiber.Ctx) error {
		item_ := &models.Question{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.QuestionController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/question", func(c *fiber.Ctx) error {
		item_ := &models.Question{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.QuestionController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/question/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.Question{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.QuestionController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Get("/question/find/monthandday/:monthandday", func(c *fiber.Ctx) error {
		monthandday_ := c.Params("monthandday")
		var controller rest.QuestionController
		controller.Init(c)
		controller.FindByMonthAndDay(monthandday_)
		controller.Close()
		return c.JSON(controller.Result)
	})

}