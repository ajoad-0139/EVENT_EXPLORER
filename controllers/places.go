package controllers

import (
	"event-explorer/requests"
	"event-explorer/utils"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type PlacesController struct {
	beego.Controller
}

func (p *PlacesController) AutoCompletedPlaces() {
	input := strings.TrimSpace(p.GetString("input"))
	token := p.GetString("sessionToken")

	// we need to check if the request is coming from the same app or not first

	// check input size
	if len(input) < 2 || len(input) > 100 {
		utils.JsonError(&p.Controller, 400, "input size must not less than 2 or greater than 100")
		return
	}
	if token == "" {
		utils.JsonError(&p.Controller, 400, "session token is not present")
		return
	}
	suggestions, statusCode, err := requests.GetAutoCompletePlaces(p.Ctx.Request.Context(), input, token)
	if err != nil {
		utils.JsonError(&p.Controller, 500, err)
		return
	}
	utils.JsonSuccess(&p.Controller, statusCode, suggestions)
}

func (p *PlacesController) GetAPlaceById() {
	placeId := strings.TrimSpace(p.Ctx.Input.Param(":placeId"))
	token := p.GetString("sessionToken")

	if placeId == "" || len(placeId) > 200 {
		utils.JsonError(&p.Controller, 400, "place id is missing or invalid")
		return
	}
	if token == "" {
		utils.JsonError(&p.Controller, 400, "session token is not present")
		return
	}

	city, statusCode, err := requests.GetPlaceById(p.Ctx.Request.Context(), placeId, token)
	if err != nil {
		utils.JsonError(&p.Controller, statusCode, err.Error())
		return
	}
	utils.JsonSuccess(&p.Controller, statusCode, city)
}
