package controllers

import (
	"event-explorer/requests"
	"event-explorer/utils"

	"github.com/beego/beego/v2/core/logs"
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
	limit := "6"
	cacheKey := requests.ListKey(city, countryCode, limit)

	//get cached data
	cachedData, err := requests.ListCache.Get(e.Ctx.Request.Context(), cacheKey)

	if err == nil && cachedData != nil {
		e.Data["events"] = cachedData
		e.TplName = "lists.tpl"
		return
	}

	res, statusCode, errMsg := requests.GetEvents(e.Ctx.Request.Context(), city, countryCode, limit)
	if errMsg != "" {
		utils.JsonError(&e.Controller, statusCode, errMsg)
		return
	}
	//set cache
	if err := requests.ListCache.Put(e.Ctx.Request.Context(), cacheKey, res, requests.ListTTL); err != nil {
		logs.Warn("list cache put failed key=%s err=%v", cacheKey, err)
	}
	e.Data["events"] = res
	e.TplName = "lists.tpl"
	// utils.JsonSuccess(&e.Controller, statusCode, res)

}

func (e *EventsController) SingleEvent() {
	eventId := e.Ctx.Input.Param(":eventId")

	if eventId == "" {
		utils.JsonError(&e.Controller, 400, "must provide eventId as a path parameter")
		return
	}
	cacheKey := requests.SingleKey(eventId)
	cachedData, err := requests.SingleCache.Get(e.Ctx.Request.Context(), cacheKey)
	if err == nil && cachedData != nil {
		e.Data["single-event"] = cachedData
		e.TplName = "single-event.tpl"
		return
	}
	res, statusCode, errMsg := requests.GetSingleEvent(e.Ctx.Request.Context(), eventId)
	if errMsg != "" {
		utils.JsonError(&e.Controller, statusCode, errMsg)
		return
	}

	if err := requests.SingleCache.Put(e.Ctx.Request.Context(), cacheKey, res, requests.SingleTTL); err != nil {
		logs.Warn("single cache put failed key=%s err=%v", cacheKey, err)
	}

	e.Data["single-event"] = res
	e.TplName = "single-event.tpl"
	// utils.JsonSuccess(&e.Controller, statusCode, res)

}
