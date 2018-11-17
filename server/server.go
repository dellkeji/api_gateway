package server

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"

	conf "apigw_golang/configure"
	storage "apigw_golang/storage"
	utils "apigw_golang/utils"
)

// channel for request
var requestInfoList = make(chan requestDestInfo, conf.GlobalConfigurations.ChannelSize)

func PrepareRequest(ctx *gin.Context, db *gorm.DB) {
	go func() {
		requestInfo := requestDestInfo{}
		// data handler
		srvCtx := serverContext{ctx: ctx}
		currCtx := srvCtx.ctx
		requestInfo.ctx = currCtx
		// api name
		apiName := currCtx.Param("api_name")
		apiInfo := gatewayInfo{}
		if err := apiInfo.MatchGateway(db, apiName); err != nil {
			requestInfo.message = fmt.Sprintf("%v", err)
			requestInfoList <- requestInfo
			return
		}
		// stage name
		stageName := currCtx.Param("stage_name")
		stage := stageInfo{}
		if err := stage.MatchStage(db, apiInfo.ID, stageName); err != nil {
			requestInfo.message = fmt.Sprintf("%v", err)
			requestInfoList <- requestInfo
			return
		}

		reqMethod := currCtx.Request.Method
		// suffix path
		suffixPath := currCtx.Param("suffix_path")

		resInfo := resourceInfo{}
		if err := resInfo.MatchReource(db, apiInfo.ID, suffixPath, reqMethod); err != nil {
			requestInfo.message = fmt.Sprintf("%v", err)
			requestInfoList <- requestInfo
			return
		}
		// render dest url
		if resInfo.dynamicPathExist {
			if err := resInfo.RenderDestURL(stage.context); err != nil {
				requestInfo.message = fmt.Sprintf("%v", err)
				requestInfoList <- requestInfo
				return
			}
		}
		// params and data
		params := currCtx.Request.URL.RawQuery
		data, err := srvCtx.dataHandler()
		if err != nil {
			requestInfo.message = fmt.Sprintf("%v", err)
			requestInfoList <- requestInfo
			return
		}
		requestInfo.method = reqMethod
		requestInfo.destURL = resInfo.destURL
		requestInfo.params = params
		requestInfo.data = data
		requestInfo.timeout = resInfo.timeout
		// input channel
		requestInfoList <- requestInfo
	}()
}

// ProxyHandler : proxy for resource
func ProxyHandler(ctx *gin.Context) {
	// panic throw upper
	// note: test the performance
	defer func() {
		if err := recover(); err != nil {
			utils.ExceptionRecovery(ctx)
			return
		}
	}()
	var err error

	db := storage.GetDBSession().DB
	db.SingularTable(true)

	// filter api
	PrepareRequest(ctx, db)

	requestInfo := requestDestInfo{}
	requestInfo = <-requestInfoList
	// error response
	if requestInfo.message != "" {
		utils.FailedResponse(requestInfo.ctx, requestInfo.message)
		return
	}
	// request dest url
	statusCode, respData, err := utils.CustomHTTPRequest(requestInfo.method, requestInfo.destURL, requestInfo.params, requestInfo.data, requestInfo.timeout)
	if err != nil {
		utils.FailedResponse(ctx, fmt.Sprintf("%v", err))
		return
	}
	utils.RegisteredServiceResponse(ctx, statusCode, respData)
	return
}
