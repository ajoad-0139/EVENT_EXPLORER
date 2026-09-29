package routers

import (
	"event-explorer/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	api := beego.NewNamespace("/api",
		beego.NSNamespace("/locations",
			beego.NSRouter("/autocomplete", &controllers.PlacesController{}, "get:AutoCompletedPlaces"),
			beego.NSRouter("/:placeId", &controllers.PlacesController{}, "get:GetAPlaceById"),
		),
	)
	beego.AddNamespace(api)
}
