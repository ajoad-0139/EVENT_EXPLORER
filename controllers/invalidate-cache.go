package controllers

import (
	"event-explorer/requests"
	"event-explorer/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type CacheController struct {
	beego.Controller
}

// DELETE /api/cache/events?city=Toronto&countryCode=CA
func (c *CacheController) InvalidateEventList() {
	city := c.GetString("city")
	countryCode := c.GetString("countryCode")

	if city == "" {
		utils.JsonError(&c.Controller, 400, "must provide city as a query parameter")
		return
	}
	if countryCode == "" {
		utils.JsonError(&c.Controller, 400, "must provide countryCode as a query parameter")
		return
	}
	limit := "6"
	cacheKey := requests.ListKey(city, countryCode, limit)

	exists, _ := requests.ListCache.IsExist(c.Ctx.Request.Context(), cacheKey)
	if !exists {
		utils.JsonSuccess(&c.Controller, 200, map[string]interface{}{"key": cacheKey, "deleted": false})
		return
	}

	if err := requests.ListCache.Delete(c.Ctx.Request.Context(), cacheKey); err != nil {
		logs.Warn("list cache delete failed key=%s err=%v", cacheKey, err)
		utils.JsonError(&c.Controller, 500, "can not delete list cache")
		return
	}

	logs.Info("list cache invalidated key=%s", cacheKey)
	utils.JsonSuccess(&c.Controller, 200, map[string]interface{}{"key": cacheKey, "deleted": true})
}

// DELETE /api/cache/events/:eventId
func (c *CacheController) InvalidateSingleEvent() {
	eventId := c.Ctx.Input.Param(":eventId")

	if eventId == "" {
		utils.JsonError(&c.Controller, 400, "must provide eventId as a path parameter")
		return
	}
	cacheKey := requests.SingleKey(eventId)

	exists, _ := requests.SingleCache.IsExist(c.Ctx.Request.Context(), cacheKey)
	if !exists {
		utils.JsonSuccess(&c.Controller, 200, map[string]interface{}{"key": cacheKey, "deleted": false})
		return
	}

	if err := requests.SingleCache.Delete(c.Ctx.Request.Context(), cacheKey); err != nil {
		logs.Warn("single cache delete failed key=%s err=%v", cacheKey, err)
		utils.JsonError(&c.Controller, 500, "can not delete single event cache")
		return
	}

	logs.Info("single cache invalidated key=%s", cacheKey)
	utils.JsonSuccess(&c.Controller, 200, map[string]interface{}{"key": cacheKey, "deleted": true})
}

// DELETE /api/cache/all
func (c *CacheController) InvalidateAll() {
	if err := requests.ListCache.ClearAll(c.Ctx.Request.Context()); err != nil {
		logs.Warn("list cache clear failed err=%v", err)
		utils.JsonError(&c.Controller, 500, "can not clear list cache")
		return
	}
	if err := requests.SingleCache.ClearAll(c.Ctx.Request.Context()); err != nil {
		logs.Warn("single cache clear failed err=%v", err)
		utils.JsonError(&c.Controller, 500, "can not clear single event cache")
		return
	}

	logs.Info("all event caches cleared")
	utils.JsonSuccess(&c.Controller, 200, map[string]interface{}{"cleared": true})
}
