package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"knowledge-base-service/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	TYPE_GITHUB             = 1
	GITHUB_CLIENT_ID        = "623037fcf1a6cb4ad6d8"
	GITHUB_CLIENT_SECRET    = "7ccd7c57dce15c44deee8760f275085afe567708"
	GITHUB_ACCESS_TOKEN_URL = "https://github.com/login/oauth/access_token"
)

func (e *User) GetProfile(ctx *gin.Context) {

}

func (e *User) UpdateProfile(ctx *gin.Context) {

}

func (e *User) SignUp(ctx *gin.Context) {

}

func (e *User) SignIn(ctx *gin.Context) {
	var payload SignInPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		tools.RespFail(ctx, 1, "参数错误:"+err.Error(), nil)
		return
	}
	if payload.Type == TYPE_GITHUB {
		tokenResp, err := getGitHubToken(payload.Code)
		if err != nil {
			tools.RespFail(ctx, 1, err.Error(), nil)
			return
		}
		tools.RespSuccess(ctx, tokenResp)
		return
	}
	tools.RespFail(ctx, 1, "未知登录类型", nil)
}

func getGitHubToken(code string) (GitHubTokenSuccessResp, error) {
	params := GitHubTokenPayload{
		ClientID:     GITHUB_CLIENT_ID,
		ClientSecret: GITHUB_CLIENT_SECRET,
		Code:         code,
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	req, err := http.NewRequest("POST", GITHUB_ACCESS_TOKEN_URL, bytes.NewBuffer(jsonParams))
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	fmt.Println("body", string(body))
	tokenResp := GitHubTokenResp{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return GitHubTokenSuccessResp{}, err
	}
	successResp := GitHubTokenSuccessResp{
		AccessToken: tokenResp.AccessToken,
		Scope:       tokenResp.Scope,
		TokenType:   tokenResp.TokenType,
	}
	if len(tokenResp.Error) > 0 {
		return GitHubTokenSuccessResp{}, errors.New(tokenResp.ErrorDescription)
	}
	return successResp, nil
}
