package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	conf "apigw_golang/configure"
	"apigw_golang/storage"
)

// Router : desc router
type Router struct {
	router *gin.Engine
}

// Init :
func (r *Router) Init(handle gin.HandlerFunc) error {
	r.router = gin.Default()
	// register method
	r.methodRegister(handle)
	hostPort := fmt.Sprintf("%v:%v", conf.GlobalConfigurations.SvrConf.Host, conf.GlobalConfigurations.SvrConf.Port)

	// connect mysql
	dbSession := storage.GetDBSession()
	defer dbSession.Close()

	return r.router.Run(hostPort)
}

func (r *Router) methodRegister(handle gin.HandlerFunc) {
	// set router group, dynamic router
	routerGroup := r.router.Group("/api/:api_name/:stage_name/*suffix_path")
	routerGroup.GET("", handle)
	routerGroup.POST("", handle)
	routerGroup.PUT("", handle)
	routerGroup.DELETE("", handle)
	routerGroup.HEAD("", handle)
}
