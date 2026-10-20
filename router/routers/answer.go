package routers

import (
	"strconv"

	"tomelater/global/log"

	"tomelater/controllers/rest"

	"tomelater/models"
	"tomelater/models/answer"

	"github.com/gofiber/fiber/v2"
)

// ---- 답변 소유권/공개 범위 ----
//
// 예전엔 로그인한 아무나 남의 uid 로 답변을 만들고, 남의 답변을 수정·삭제하고, 비공개 답변을
// 목록·조회로 읽을 수 있었다. 이제는
//   - 작성: 작성자(uid)를 토큰의 사용자로 고정
//   - 수정·삭제: 작성자 본인만(아니면 403)
//   - 조회: 공개 답변 + 본인 답변만, 남의 답변에 딸린 작성자 이메일/소셜 ID 는 지움

// canViewAnswer reports whether viewer may read a.
func canViewAnswer(a *models.Answer, viewer int64) bool {
	return a != nil && (a.Ispublic == answer.IspublicPublic || (viewer != 0 && a.Uid == viewer))
}

// stripAnswerAuthor blanks the embedded author's e-mail / social id unless the
// viewer is the author.
func stripAnswerAuthor(a *models.Answer, viewer int64) {
	if a == nil || a.Extra == nil || a.Uid == viewer {
		return
	}
	if u, ok := a.Extra["user"].(models.User); ok {
		u.Email, u.Socialid = "", ""
		a.Extra["user"] = u
	}
}

func visibleAnswers(items []models.Answer, viewer int64) []models.Answer {
	out := make([]models.Answer, 0, len(items))
	for i := range items {
		if canViewAnswer(&items[i], viewer) {
			stripAnswerAuthor(&items[i], viewer)
			out = append(out, items[i])
		}
	}
	return out
}

// filterAnswerResult applies the visibility rule to a rest.AnswerController
// result ("item" / "items").
func filterAnswerResult(result map[string]interface{}, viewer int64) {
	switch t := result["item"].(type) {
	case *models.Answer:
		if t != nil {
			if !canViewAnswer(t, viewer) {
				result["item"] = nil
			} else {
				stripAnswerAuthor(t, viewer)
			}
		}
	}
	switch t := result["items"].(type) {
	case []models.Answer:
		result["items"] = visibleAnswers(t, viewer)
	case *[]models.Answer:
		if t != nil {
			result["items"] = visibleAnswers(*t, viewer)
		}
	}
}

// ownsAnswer reports whether the answer exists and belongs to uid.
func ownsAnswer(id int64, uid int64) bool {
	if id == 0 || uid == 0 {
		return false
	}
	manager := models.NewAnswerManager(nil)
	defer manager.Close()
	a := manager.Get(id)
	return a != nil && a.Uid == uid
}

func answerForbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "error", "message": "forbidden"})
}

func answerUnauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "error", "message": "unauthorized"})
}

// SetupAnswerRoutes sets up routes for answer domain
func SetupAnswerRoutes(group fiber.Router) {

	group.Get("/answer", func(c *fiber.Ctx) error {
		page_, _ := strconv.Atoi(c.Query("page"))
		pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
		var controller rest.AnswerController
		controller.Init(c)
		controller.Index(page_, pagesize_)
		controller.Close()
		// total 은 필터 전 개수라 남의 비공개 답변 수만큼 클 수 있다(내용은 노출되지 않음).
		filterAnswerResult(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

	group.Get("/answer/:id", func(c *fiber.Ctx) error {
		id_, _ := strconv.ParseInt(c.Params("id"), 10, 64)
		var controller rest.AnswerController
		controller.Init(c)
		controller.Read(id_)
		controller.Close()
		filterAnswerResult(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

	group.Post("/answer", func(c *fiber.Ctx) error {
		me := viewerUserID(c)
		if me == 0 {
			return answerUnauthorized(c)
		}
		item_ := &models.Answer{}
		err := c.BodyParser(item_)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		// 작성자는 토큰의 사용자로 고정한다(남의 uid 로 답변을 만들지 못하게).
		item_.Uid = me
		var controller rest.AnswerController
		controller.Init(c)
		controller.Insert(item_)
		controller.Close()
		return c.JSON(controller.Result)
	})

	group.Post("/answer/batch", func(c *fiber.Ctx) error {
		me := viewerUserID(c)
		if me == 0 {
			return answerUnauthorized(c)
		}
		var items_ *[]models.Answer
		items__ref := &items_
		err := c.BodyParser(items__ref)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if items_ != nil {
			for i := range *items_ {
				(*items_)[i].Uid = me
			}
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
		me := viewerUserID(c)
		if !ownsAnswer(item_.Id, me) {
			return answerForbidden(c)
		}
		// uid 는 바꿀 수 없다(답변 소유권 이전 방지).
		item_.Uid = me
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
		if !ownsAnswer(item_.Id, viewerUserID(c)) {
			return answerForbidden(c)
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
		me := viewerUserID(c)
		for _, it := range *item_ {
			if !ownsAnswer(it.Id, me) {
				return answerForbidden(c)
			}
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
		filterAnswerResult(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

	group.Get("/answer/find/qid/:qid", func(c *fiber.Ctx) error {
		qid_, _ := strconv.ParseInt(c.Params("qid"), 10, 64)
		var controller rest.AnswerController
		controller.Init(c)
		controller.FindByQid(qid_)
		controller.Close()
		filterAnswerResult(controller.Result, viewerUserID(c))
		return c.JSON(controller.Result)
	})

}
