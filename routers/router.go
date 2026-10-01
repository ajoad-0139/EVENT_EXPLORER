package routers

import (
	"event-explorer/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//frontend routes
	beego.Router("/", &controllers.MainController{})
	beego.Router("/events", &controllers.EventsController{}, "get:ListEvents")
	beego.Router("/events/:eventId", &controllers.EventsController{}, "get:SingleEvent")
	beego.Router("/redirect/:eventId", &controllers.EventsController{}, "get:RedirectToTickets")

	//api proxies for serving data to browser
	api := beego.NewNamespace("/api",
		beego.NSNamespace("/locations",
			beego.NSRouter("/autocomplete", &controllers.PlacesController{}, "get:AutoCompletedPlaces"),
			beego.NSRouter("/:placeId", &controllers.PlacesController{}, "get:GetAPlaceById"),
		),
	)
	beego.AddNamespace(api)
}
