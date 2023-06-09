package wallpaper

import (
	"encoding/json"
	"io/ioutil"
	"knowledge-base-service/tools"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	WALLHAVEN_API = "https://wallhaven.cc/api/v1"
)

func (e *Wallpaper) Search(ctx *gin.Context) {
	var query SearchQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	resp, err := http.Get(WALLHAVEN_API + "/search?" + ctx.Request.URL.RawQuery)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	var searchResult WallhavenSearchResp
	if err = json.Unmarshal(body, &searchResult); err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, searchResult.Data)
}

func (e *Wallpaper) GetInfo(ctx *gin.Context) {
	var query GetInfoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	decodedURL, err := url.QueryUnescape(query.URL)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	url := strings.Replace(decodedURL, "https://wallhaven.cc", WALLHAVEN_API, 1)
	resp, err := http.Get(url)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	var infoResult WallhavenInfoResp
	if err := json.Unmarshal(body, &infoResult); err != nil {
		tools.RespFail(ctx, 1, err.Error(), nil)
		return
	}
	tools.RespSuccess(ctx, infoResult.Data)
}
