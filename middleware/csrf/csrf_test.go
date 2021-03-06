package csrf

import (
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

func Test_CSRF(t *testing.T) {
	app := fiber.New()

	app.Use(New())

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	h := app.Handler()
	ctx := &fasthttp.RequestCtx{}

	// Generate CSRF token
	ctx.Request.Header.SetMethod("GET")
	h(ctx)
	utils.AssertEqual(t, true, strings.Contains(string(ctx.Response.Header.Peek(fiber.HeaderSetCookie)), "_csrf"))

	// Without CSRF cookie
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod("POST")
	h(ctx)
	utils.AssertEqual(t, 403, ctx.Response.StatusCode())

	// Empty/invalid CSRF token
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.Set("X-CSRF-Token", "johndoe")
	h(ctx)
	utils.AssertEqual(t, 403, ctx.Response.StatusCode())

	// Valid CSRF token
	token := utils.UUID()
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.Set(fiber.HeaderCookie, "_csrf="+token)
	ctx.Request.Header.Set("X-CSRF-Token", token)
	h(ctx)
	utils.AssertEqual(t, 200, ctx.Response.StatusCode())
}
