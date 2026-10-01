package utils

import (
	"event-explorer/models"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// RenderErrorMsg renders error.tpl with a message written for the visitor
func RenderErrorMsg(c *beego.Controller, status int, title string, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["error"] = models.ErrorPage{
		Status:  status,
		Title:   title,
		Message: message,
	}
	c.TplName = "error.tpl"

	if err := c.Render(); err != nil {
		logs.Error("error page render failed err=%v", err)
		c.Ctx.WriteString(message)
	}
}

func RenderError(c *beego.Controller, status int, errMsg string) {
	logs.Error("request failed status=%d err=%s path=%s", status, errMsg, c.Ctx.Request.URL.Path)

	switch {
	case status == http.StatusBadRequest:
		RenderErrorMsg(c, status, "Something is missing", "The link looks incomplete. Please go back and start a new search.")
	case status == http.StatusNotFound:
		RenderErrorMsg(c, status, "Event not found", "We could not find the event you were looking for. It may have been removed.")
	case status == http.StatusTooManyRequests:
		RenderErrorMsg(c, status, "Too many requests", "Please wait a moment and try again.")
	case status == http.StatusGatewayTimeout || status == http.StatusRequestTimeout:
		RenderErrorMsg(c, http.StatusGatewayTimeout, "Taking too long", "The event provider is slow to respond. Please try again in a moment.")
	default:
		RenderErrorMsg(c, http.StatusBadGateway, "Something went wrong", "We could not load events right now. Please try again in a moment.")
	}
}
