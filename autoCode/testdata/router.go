package initialize

import "github.com/gin-gonic/gin"

func Routers() *gin.Engine {
	var Router = gin.Default()
	g := Router.Group("")
	{
		router.InitApiRouter(g)
		// Code generated by XXX Begin; DO NOT EDIT.
		// Code generated by XXX End; DO NOT EDIT.
	}
	return Router
}
