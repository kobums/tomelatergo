package routers

import (

	"encoding/json"

	"strconv"

	"tomelater/global/log"


	"tomelater/controllers/rest"


	"tomelater/models/user"

	"tomelater/models"
	"github.com/gofiber/fiber/v2"
)

// viewerUserID returns the id of the authenticated caller (set by JwtAuthRequired), or 0.
func viewerUserID(c *fiber.Ctx) int64 {
	if u, ok := c.Locals("user").(*models.User); ok && u != nil {
		return u.Id
	}
	return 0
}

// hideOthersEmail blanks Email / Socialid of every user in result that is not
// the caller. 예전엔 로그인한 아무나 /api/user 목록·조회로 전체 회원의 이메일을 가져갈 수 있었다.
func hideOthersEmail(result map[string]interface{}, me int64) {
	strip := func(u *models.User) {
		if u != nil && u.Id != me {
			u.Email, u.Socialid = "", ""
		}
	}
	for _, v := range result {
		switch t := v.(type) {
		case *models.User:
			strip(t)
		case []models.User:
			for i := range t {
				strip(&t[i])
			}
		case *[]models.User:
			if t != nil {
				for i := range *t {
					strip(&(*t)[i])
				}
			}
		}
	}
	// 값 타입(models.User)은 복사본이므로 고친 뒤 다시 넣는다.
	for k, v := range result {
		if t, ok := v.(models.User); ok {
			strip(&t)
			result[k] = t
		}
	}
}

// SetupUserRoutes sets up routes for user domain
func SetupUserRoutes(group fiber.Router) {

	group.Get("/user", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.UserController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		hideOthersEmail(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

	group.Get("/user/:id", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.UserController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		hideOthersEmail(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

	group.Post("/user", func(c *fiber.Ctx) error {
		item_ := &models.User{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.UserController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/user/batch", func(c *fiber.Ctx) error {
		var items_ *[]models.User
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.UserController
		controller.Init(c)
		controller.Insertbatch(items_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/user/count", func(c *fiber.Ctx) error {

		var controller rest.UserController
		controller.Init(c)
		controller.Count()
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Put("/user", func(c *fiber.Ctx) error {
		item_ := &models.User{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.UserController
		controller.Init(c)
		controller.Update(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/user", func(c *fiber.Ctx) error {
		item_ := &models.User{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.UserController
		controller.Init(c)
		controller.Delete(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Delete("/user/batch", func(c *fiber.Ctx) error {
		item_ := &[]models.User{}
		err := c.BodyParser(item_)
		if err != nil {
		    log.Error().Msg(err.Error())
		}
		var controller rest.UserController
		controller.Init(c)
		controller.Deletebatch(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	// /user/find/email, /user/find/socialid 는 제거했다(다른 회원의 이메일·소셜 계정 조회/존재 확인).

	group.Get("/user/find/nickname/:nickname", func(c *fiber.Ctx) error {
		nickname_ := c.Params("nickname")
		var controller rest.UserController
		controller.Init(c)
		controller.FindByNickname(nickname_)
		controller.Close()
		hideOthersEmail(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

}