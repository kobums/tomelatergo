package router

import (
	"strconv"
	"strings"

	"tomelater/router/routers"

	"github.com/gofiber/fiber/v2"
)

func getArrayCommal(name string) []int64 {
	values := strings.Split(name, ",")

	var items []int64
	for _, item := range values {
        n, _ := strconv.ParseInt(item, 10, 64)
		items = append(items, n)
	}

	return items
}

func getArrayCommai(name string) []int {
	values := strings.Split(name, ",")

	var items []int
	for _, item := range values {
        n, _ := strconv.Atoi(item)
		items = append(items, n)
	}

	return items
}

func SetRouter(r *fiber.App) {

	// Deprecated: 비밀번호가 URL 쿼리스트링에 실려 프록시·액세스 로그에 남는다.
	// 새 클라이언트는 POST /api/auth/login 을 쓴다. 배포된 앱 호환용으로만 남겨 둔다
	// (요청 로그에서는 제외: services/http.go 의 SkipURIs).
	r.Get("/api/jwt", func(c *fiber.Ctx) error {
		loginid := c.Query("loginid")
		passwd := c.Query("passwd")
		return c.JSON(JwtAuth(c, c.Query("type"), c.Query("token"), loginid, passwd))
	})

	// 로그인: 자격 증명을 JSON 본문으로 받는다.
	//   {"loginid": "...", "passwd": "..."}  또는  {"type": "apple|google|kakao|naver", "token": "..."}
	r.Post("/api/auth/login", func(c *fiber.Ctx) error {
		var body struct {
			Loginid string `json:"loginid"`
			Email   string `json:"email"`
			Passwd  string `json:"passwd"`
			Type    string `json:"type"`
			Token   string `json:"token"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "error", "message": "invalid body"})
		}
		loginid := body.Loginid
		if loginid == "" {
			loginid = body.Email
		}
		return c.JSON(JwtAuth(c, body.Type, body.Token, loginid, body.Passwd))
	})

	apiGroup := r.Group("/api")

	apiGroup.Use(JwtAuthRequired)


	// Setup domain-specific routes
	routers.SetupMembershipRoutes(apiGroup)
	routers.SetupMembershipusageRoutes(apiGroup)
	routers.SetupHealthRoutes(apiGroup)
	routers.SetupSettingRoutes(apiGroup)
	routers.SetuptomelatertrainerRoutes(apiGroup)
	routers.SetupStopRoutes(apiGroup)
	routers.SetupRockergroupRoutes(apiGroup)
	routers.SetupTrainermemberRoutes(apiGroup)
	routers.SetupTokenRoutes(apiGroup)
	routers.SetupPaymenttypeRoutes(apiGroup)
	routers.SetupUsehealthusageRoutes(apiGroup)
	routers.SetupInquiryRoutes(apiGroup)
	routers.SetupHealthcategoryRoutes(apiGroup)
	routers.SetupOrderRoutes(apiGroup)
	routers.SetupNoticeRoutes(apiGroup)
	routers.SetupTermRoutes(apiGroup)
	routers.SetuptomelaterRoutes(apiGroup)
	routers.SetupDaytypeRoutes(apiGroup)
	routers.SetupLoginlogRoutes(apiGroup)
	routers.SetupPtreservationRoutes(apiGroup)
	routers.SetupAppversionRoutes(apiGroup)
	routers.SetupQrcodeRoutes(apiGroup)
	routers.SetupSystemlogRoutes(apiGroup)
	routers.SetupDiscountRoutes(apiGroup)
	routers.SetupRockerRoutes(apiGroup)
	routers.SetupRoleRoutes(apiGroup)
	routers.SetupPaymentRoutes(apiGroup)
	routers.SetupRockerusageRoutes(apiGroup)
	routers.SetupUsehealthRoutes(apiGroup)
	routers.SetupUserRoutes(apiGroup)
	routers.SetupPaymentformRoutes(apiGroup)
	routers.SetupIpblockRoutes(apiGroup)
	routers.SetupAlarmRoutes(apiGroup)
	routers.SetupMemberqrRoutes(apiGroup)
	routers.SetupWorkoutlogRoutes(apiGroup)
	routers.SetupMemberbodyRoutes(apiGroup)
	routers.SetupPushtokenRoutes(apiGroup)
	routers.SetupAttendanceRoutes(apiGroup)
}