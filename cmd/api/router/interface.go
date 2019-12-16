package router

// Interface receives a new Router through method
type Router interface {
	StartWebServerHTTP()
	SetupRoutes()
}
