package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-zoox/fetch"
	uicommon "github.com/oustn/qtg/ui/common"
)

const deviceId = "66f6e3b560ad8876e52e6e67ee535c5c"
const UA = `QingTing-iOS/10.4.2.3 com.Qting.QTTour Mozilla/5.0 (iPhone; CPU iPhone OS 16_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`

type QingTingApi struct {
	RefreshToken string
	AccessToken  string
	ExpiresIn    int
	QingTingId   string
	User         uicommon.UserInfo
}

type ResponseData struct {
	ErrorNo  int         `json:"errorno"`
	ErrorMsg string      `json:"errormsg"`
	Data     interface{} `json:"data"`
}

func (api *QingTingApi) get(url string, config *fetch.Config, result interface{}) error {
	config.Headers["User-Agent"] = UA
	response, err := fetch.Get(url, config)

	if err != nil {
		return err
	}

	data, err := response.JSON()
	if err != nil {
		return err
	}

	respData := &ResponseData{
		Data: result,
	}
	err = json.Unmarshal([]byte(data), respData)
	if err != nil {
		return err
	}

	if respData.ErrorNo != 0 {
		return fmt.Errorf("error %d: %s", respData.ErrorNo, respData.ErrorMsg)
	}

	return nil
}

func (api *QingTingApi) post(url string, config *fetch.Config, result interface{}) error {
	config.Headers["User-Agent"] = UA
	response, err := fetch.Post(url, config)
	if err != nil {
		return err
	}

	data, err := response.JSON()
	if err != nil {
		return err
	}

	respData := &ResponseData{
		Data: result,
	}
	err = json.Unmarshal([]byte(data), respData)
	if err != nil {
		return err
	}

	if respData.ErrorNo != 0 {
		return fmt.Errorf("error %d: %s", respData.ErrorNo, respData.ErrorMsg)
	}

	return nil
}

func InitQingTingApi(options ...string) *QingTingApi {
	api := &QingTingApi{
		RefreshToken: options[0],
		QingTingId:   options[1],
	}
	return api
}

func (api *QingTingApi) Auth(options ...string) error {
	type Auth struct {
		QingtingID   string `json:"qingting_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	var (
		auth  Auth
		token string
		id    string
	)
	if len(options) > 0 {
		token = options[0]
	}
	if token == "" {
		token = api.RefreshToken
	}
	if len(options) > 1 {
		id = options[1]
	}
	if id == "" {
		id = api.QingTingId
	}

	if token == "" {
		return fmt.Errorf("RefreshToken 为空")
	}
	if id == "" {
		return fmt.Errorf("QingTingId 为空")
	}

	err := api.post("https://user.qtfm.cn/u2/api/v4/auth", &fetch.Config{
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent":   UA,
		},
		Body: map[string]string{
			"refresh_token": token,
			"qingting_id":   id,
			"device_id":     deviceId,
			"grant_type":    "refresh_token",
		},
	}, &auth)

	if err != nil {
		return err
	}
	api.RefreshToken = auth.RefreshToken
	api.QingTingId = auth.QingtingID
	api.AccessToken = auth.AccessToken
	api.ExpiresIn = auth.ExpiresIn
	return api.FetchUserInfo()
}

func (api *QingTingApi) FetchUserInfo() error {
	if api.RefreshToken == "" {
		return fmt.Errorf("未授权")
	}

	var userInfo uicommon.UserInfo
	url := "https://user.qtfm.cn/u2/api/v5/user/" + api.QingTingId + "?device_id=" + deviceId + "&mode=vital&qingting_id=" + api.QingTingId + "&access_token=" + api.AccessToken
	err := api.get(url, &fetch.Config{
		Headers: map[string]string{
			"Authorization": "Bearer " + api.AccessToken,
			"User-Agent":    UA,
		},
	}, &userInfo)

	if err != nil {
		return err
	}
	api.User = userInfo
	return nil
}
