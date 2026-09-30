package controllers

import (
	"event-explorer/requests"
	"event-explorer/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type EventsController struct {
	beego.Controller
}

func (e *EventsController) ListEvents() {
	city := e.GetString("city")
	countryCode := e.GetString("countryCode")

	if city == "" {
		utils.JsonError(&e.Controller, 400, "must provide city as a query parameter")
		return
	}
	if countryCode == "" {
		utils.JsonError(&e.Controller, 400, "must provide countryCode as a query parameter")
		return
	}

	res, statusCode, errMsg := requests.GetEvents(e.Ctx.Request.Context(), city, countryCode, "6")
	if errMsg != "" {
		utils.JsonError(&e.Controller, statusCode, errMsg)
		return
	}

	e.Data["events"] = res
	e.TplName = "lists.tpl"
	// utils.JsonSuccess(&e.Controller, statusCode, res)

}

func (e*EventsController) SingleEvent() {
	
}
