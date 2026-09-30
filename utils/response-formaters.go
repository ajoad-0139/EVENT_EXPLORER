package utils

import (
	"time"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.AddFuncMap("formatEventDate", FormatEventDate)
}

// reusable response formaters

// error response
func JsonError(c *web.Controller, status int, err interface{}) {
	c.Data["json"] = map[string]interface{}{
		"Error": err,
	}
	c.Ctx.Output.SetStatus(status)
	c.ServeJSON()
}

// success response
func JsonSuccess(c *web.Controller, status int, data interface{}) {
	c.Data["json"] = data
	c.Ctx.Output.SetStatus(status)
	c.ServeJSON()
}

func FormatEventDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("Mon, 02 Jan 2006")
}
