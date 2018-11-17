package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/jinzhu/gorm"

	"apigw_golang/storage/models"
	utils "apigw_golang/utils"
)

const (
	matchRule           = `\{([A-Za-z0-9_-]+)\}`
	subReplaceRule      = `(?P<$1>[^/]+)`
	destPathRule        = `\{[stageVariables\.]*?([A-Za-z0-9_-]+)\}`
	destPathReplaceRule = `{{$1}}`
)

// dataHandler : get data
func (ctx *serverContext) dataHandler() (data string, err error) {
	byteData, err := utils.GetRawData(ctx.ctx)
	return string(byteData), err
}

// MatchGateway : match the gateway(the group of resource)
func (apiInfo *gatewayInfo) MatchGateway(db *gorm.DB, apiName string) error {
	var apiDetail models.API
	if err := db.Where(
		"api_name = ?", apiName,
	).First(&apiDetail).Error; err != nil {
		return err
	}
	apiInfo.apiName = apiDetail.APIName
	apiInfo.apiSlugName = apiDetail.APINameSlug
	apiInfo.ID = apiDetail.ID
	return nil
}

// MatchStage : match the stage info
func (stage *stageInfo) MatchStage(db *gorm.DB, apiID int, name string) error {
	var stageDetail models.Stage
	if err := db.Where(
		"api_id = ? AND stage_name = ?", apiID, name,
	).First(&stageDetail).Error; err != nil {
		return err
	}
	// base info for stage
	stage.apiID = apiID
	stage.ID = stageDetail.ID
	stage.name = name
	//
	variables := make(pongo2.Context)
	if err := json.Unmarshal(
		[]byte(stageDetail.StageVariables), &variables,
	); err != nil {
		return err
	}
	stage.context = variables
	return nil
}

// MatchReource : match the registered resource
func (resInfo *resourceInfo) MatchReource(db *gorm.DB, apiID int, path string, method string) error {
	var resDetail models.Resource
	var resDetailList []models.Resource
	if err := db.Where(
		"api_id = ? AND registed_http_method = ?", apiID, method,
	).Find(&resDetailList).Error; err != nil {
		return err
	}

	pathExist := false
	var dynamicPathList []models.Resource
	// match dynamicPathList
	for _, info := range resDetailList {
		if strings.Contains(info.Path, "{") && strings.Contains(info.Path, "}") {
			dynamicPathList = append(dynamicPathList, info)
		} else {
			if info.Path == path {
				resDetail = info
				pathExist = true
				break
			}
		}
	}

	if !pathExist {
		dynamicPathMatch(resInfo, dynamicPathList, path, &resDetail, &pathExist)
	}
	if !pathExist {
		return fmt.Errorf("not find path")
	}

	resInfo.apiID = apiID
	resInfo.destMethod = resDetail.DestMethod
	resInfo.destURL = resDetail.DestURL
	resInfo.ID = resDetail.ID
	resInfo.path = resDetail.Path
	resInfo.registeredMethod = resDetail.RegisteredMethod
	resInfo.timeout = resDetail.Timeout
	return nil
}

// dynaminc path
func dynamicPathMatch(resInfo *resourceInfo, dynamicPathList []models.Resource, path string, resDetail *models.Resource, pathExist *bool) {
	resInfo.context = make(pongo2.Context)
	reg := regexp.MustCompile(matchRule)
	for _, info := range dynamicPathList {
		renderPath := reg.ReplaceAllString(info.Path, subReplaceRule)
		renderRegPath := regexp.MustCompile("^" + renderPath + "$")
		if renderRegPath.Match([]byte(path)) {
			subMatch := renderRegPath.FindStringSubmatch(path)
			nameMatch := renderRegPath.SubexpNames()
			for i := 1; i < len(nameMatch); i++ {
				resInfo.context[nameMatch[i]] = subMatch[i]
				resInfo.dynamicPathExist = true
			}
			*resDetail = info
			*pathExist = true
			break
		}
	}
}

// RenderDestURL : render registered resource dest url
func (resInfo *resourceInfo) RenderDestURL(stageContext pongo2.Context) error {
	tempURL := resInfo.destURL
	// replace path
	reg := regexp.MustCompile(destPathRule)
	tempURL = reg.ReplaceAllString(tempURL, destPathReplaceRule)
	// replace domain
	domainInfo, err := refineStageVariable(stageContext)
	if err != nil {
		return err
	}

	mergeContext(domainInfo, resInfo.context)
	// render
	resInfo.destURL, err = renderByPongo(tempURL, domainInfo)
	if err != nil {
		return err
	}
	return nil
}

// refineStageVariable : default use array[0]
// TODO: support multi domain
func refineStageVariable(stageContenx pongo2.Context) (refineInfo pongo2.Context, err error) {
	templInfo := make(pongo2.Context)
	for info, value := range stageContenx {
		val, ok := (value).(string)
		if !ok {
			return nil, fmt.Errorf("value type error")
		}
		list, err := utils.String2List(val)
		if err != nil {
			return nil, err
		}
		// random filter
		templInfo[info] = list[rand.Intn(len(list))]
	}
	return templInfo, nil
}

// renderByPongo : render
func renderByPongo(url string, context pongo2.Context) (str string, err error) {
	tpl, err := pongo2.FromString(url)
	// hide the raw err; replace with platform message
	if err != nil {
		return "", err
	}
	return tpl.Execute(context)
}

// mergeContext : merge two context(pongo2.context)
func mergeContext(mergedContext pongo2.Context, context pongo2.Context) {
	for info, value := range context {
		mergedContext[info] = value
	}
}
