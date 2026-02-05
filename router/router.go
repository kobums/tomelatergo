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

    r.Get("/api/jwt", func(c *fiber.Ctx) error {
		loginid := c.Query("loginid")
        passwd := c.Query("passwd")
        return c.JSON(JwtAuth(c, loginid, passwd))
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