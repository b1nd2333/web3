package bypassCaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b1nd2333/web3/pkg/common"
	"io"
	"net/http"
)

type CFCaptchaStruct struct {
	Href     string `json:"href"`
	SiteKey  string `json:"sitekey"`
	Action   string `json:"action"`
	Explicit bool   `json:"explicit"`
}

type CFCaptchaRespStruct struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Id     string `json:"id"`
	Cost   string `json:"cost"`
	Data   struct {
		Token   string `json:"token"`
		Cookies string `json:"cookies"`
	} `json:"data"`
	Extra struct {
		SecChUa                string `json:"sec-ch-ua"`
		UserAgent              string `json:"user-agent"`
		SecChUaPlatform        string `json:"sec-ch-ua-platform"`
		SecChUaArch            string `json:"sec-ch-ua-arch"`
		SecChUaBitness         string `json:"sec-ch-ua-bitness"`
		SecChUaFullVersion     string `json:"sec-ch-ua-full-version"`
		SecChUaFullVersionList string `json:"sec-ch-ua-full-version-list"`
		SecChUaMobile          string `json:"sec-ch-ua-mobile"`
		SecChUaModel           string `json:"sec-ch-ua-Model"`
		SecChUaPlatformVersion string `json:"sec-ch-ua-platform-version"`
	} `json:"extra"`
	EnMsg string `json:"en_msg"`
}

// ByPassCFCaptcha 绕过CF验证码
func ByPassCFCaptcha(proxyStr, userToken, href, siteKey, action string, explicit bool) (string, error) {
	// 创建 HTTP 客户端
	client, err := common.NewHTTPClientWithProxy(proxyStr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	cfCaptchaModel := &CFCaptchaStruct{}
	cfCaptchaModel.Href = href
	cfCaptchaModel.SiteKey = siteKey
	cfCaptchaModel.Action = action
	cfCaptchaModel.Explicit = explicit

	marshal, _ := json.Marshal(cfCaptchaModel)

	// 判断今日是否签到
	req, err := http.NewRequest("POST", "http://api.nocaptcha.io/api/wanda/cloudflare/universal", bytes.NewReader(marshal))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Token", userToken)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return "", err
	}

	cfCaptchaRespModel := &CFCaptchaRespStruct{}
	err = json.Unmarshal(body, cfCaptchaRespModel)
	if err != nil {
		return "", err
	}

	return cfCaptchaRespModel.Data.Token, nil
}
