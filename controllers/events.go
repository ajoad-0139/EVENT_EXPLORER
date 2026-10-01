package controllers

import (
	"event-explorer/models"
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
		utils.RenderErrorMsg(&e.Controller, 400, "City is missing", "Please select a city to see its events.")
		return
	}
	if countryCode == "" {
		utils.RenderErrorMsg(&e.Controller, 400, "City is missing", "Please select a city from the suggestions to see its events.")
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
		utils.RenderError(&e.Controller, statusCode, errMsg)
		return
	}
	//set cache
	if err := requests.ListCache.Put(e.Ctx.Request.Context(), cacheKey, res, requests.ListTTL); err != nil {
		logs.Warn("list cache put failed key=%s err=%v", cacheKey, err)
	}
	e.Data["events"] = res
	e.TplName = "lists.tpl"
}

func (e *EventsController) SingleEvent() {
	eventId := e.Ctx.Input.Param(":eventId")

	if eventId == "" {
		utils.RenderErrorMsg(&e.Controller, 400, "Event is missing", "The link looks incomplete. Please go back and pick an event.")
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
		utils.RenderError(&e.Controller, statusCode, errMsg)
		return
	}

	if err := requests.SingleCache.Put(e.Ctx.Request.Context(), cacheKey, res, requests.SingleTTL); err != nil {
		logs.Warn("single cache put failed key=%s err=%v", cacheKey, err)
	}

	e.Data["single-event"] = res
	e.TplName = "single-event.tpl"
}

func (e *EventsController) RedirectToTickets() {
	eventId := e.Ctx.Input.Param(":eventId")

	if eventId == "" {
		utils.RenderErrorMsg(&e.Controller, 400, "Event is missing", "The link looks incomplete. Please go back and pick an event.")
		return
	}

	var event *models.EventResponse
	cacheKey := requests.SingleKey(eventId)

	//get cached data
	cachedData, err := requests.SingleCache.Get(e.Ctx.Request.Context(), cacheKey)
	if err == nil && cachedData != nil {
		if cachedEvent, ok := cachedData.(*models.EventResponse); ok {
			event = cachedEvent
		}
	}

	if event == nil {
		res, statusCode, errMsg := requests.GetSingleEvent(e.Ctx.Request.Context(), eventId)
		if errMsg != "" {
			utils.RenderError(&e.Controller, statusCode, errMsg)
			return
		}
		//set cache
		if err := requests.SingleCache.Put(e.Ctx.Request.Context(), cacheKey, res, requests.SingleTTL); err != nil {
			logs.Warn("single cache put failed key=%s err=%v", cacheKey, err)
		}
		event = res
	}

	if event.URL == "" {
		utils.RenderErrorMsg(&e.Controller, 404, "Tickets unavailable", "Tickets are not available for this event right now.")
		return
	}

	if !utils.IsApprovedTicketURL(event.URL) {
		logs.Warn("ticket url rejected eventId=%s url=%s", eventId, event.URL)
		utils.RenderErrorMsg(&e.Controller, 502, "Ticket link blocked", "For your safety, we could not open this ticket link.")
		return
	}

	e.Redirect(event.URL, 302)
}
