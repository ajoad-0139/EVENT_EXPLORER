package utils

import (
	"github.com/beego/beego/v2/server/web"
)

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
