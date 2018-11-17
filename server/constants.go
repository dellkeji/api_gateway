package server

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

type serverContext struct {
	ctx *gin.Context
}

// group info
type gatewayInfo struct {
	ID          int
	apiName     string
	apiSlugName string
}

// resource info
type resourceInfo struct {
	ID               int
	apiID            int
	path             string
	registeredMethod string
	destMethod       string
	destURL          string
	timeout          int
	dynamicPathExist bool
	context          pongo2.Context
}

// stage info
type stageInfo struct {
	ID      int
	apiID   int
	name    string
	context pongo2.Context
}

type requestDestInfo struct {
	ctx     *gin.Context
	method  string
	destURL string
	params  string
	data    string
	timeout int
	message string
}
