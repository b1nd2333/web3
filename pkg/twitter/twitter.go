package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b1nd2333/web3/pkg/common"
	"io"
	"net/http"
	"strings"
)

// BindTwitter 绑定推特账号
func BindTwitter(uri string, xToken string, proxyStr string) (string, error) {
	uri = strings.Replace(uri, "twitter.com", "x.com", -1)

	// 获取csrfToken(ct0)
	csrfToken, err := getCsrfCode(proxyStr, uri, xToken)
	if err != nil {
		return "", err
	}

	// 获取xAuthCode
	uri1 := strings.Replace(uri, "/i/oauth2", "/i/api/2/oauth2", -1)
	code, err := getXAuthCode(proxyStr, uri1, xToken, csrfToken)
	if err != nil {
		return "", err
	}

	// 返回绑定域名
	bindUri, err := authorize(proxyStr, uri, xToken, code, csrfToken)
	if err != nil {
		return "", err
	}

	return bindUri, nil
}

// 获取csrfToken
func getCsrfCode(proxyStr string, uri string, xToken string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}

	client, _ := common.NewHTTPClientWithProxy(proxyStr)
	defer client.CloseIdleConnections()
	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36")
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s", xToken))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "ct0" {
			if cookie.Value != "" {
				return cookie.Value, nil
			}
			return "", errors.New("获取失败")
		}
	}

	return "", errors.New("获取失败")
}

type xAuthCodeRespStruct struct {
	AuthCode                          string `json:"auth_code"`
	AppName                           string `json:"app_name"`
	AppDescription                    string `json:"app_description"`
	AppUri                            string `json:"app_uri"`
	AppImageUri                       string `json:"app_image_uri"`
	OrganizationName                  string `json:"organization_name"`
	OrganizationTermsAndConditionsUri string `json:"organization_terms_and_conditions_uri"`
	OrganizationPrivacyPolicyUri      string `json:"organization_privacy_policy_uri"`
	UserCountRange                    string `json:"user_count_range"`
	Scopes                            []struct {
		Name        string `json:"name"`
		Rank        int    `json:"rank"`
		Category    string `json:"category"`
		Description string `json:"description"`
	} `json:"scopes"`
}

// 获取xAuthCode
func getXAuthCode(proxyStr string, uri string, xToken string, csrfToken string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36")
	req.Header.Set("X-Csrf-Token", csrfToken)
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s;ct0=%s", xToken, csrfToken))

	client, _ := common.NewHTTPClientWithProxy(proxyStr)
	defer client.CloseIdleConnections()

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err

	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("获取失败")
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	XAuthCodeRespModel := &xAuthCodeRespStruct{}
	json.Unmarshal(body, XAuthCodeRespModel)
	if XAuthCodeRespModel.AuthCode == "" {
		return "", errors.New("获取失败")
	}
	return XAuthCodeRespModel.AuthCode, nil
}

type authorizeRespStruct struct {
	RedirectUri string `json:"redirect_uri"`
}

// 获取callback地址
func authorize(proxyStr string, uri string, xToken string, code string, csrfToken string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("POST", "https://x.com/i/api/2/oauth2/authorize", strings.NewReader(fmt.Sprintf("approval=true&code=%s", code)))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client, _ := common.NewHTTPClientWithProxy(proxyStr)
	defer client.CloseIdleConnections()

	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36")
	req.Header.Set("X-Csrf-Token", csrfToken)
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s;ct0=%s", xToken, csrfToken))
	req.Header.Set("Referer", uri)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err

	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 401 {
			return "", errors.New("401")
		}
		return "", errors.New("获取失败")
	}
	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return "", err
	}

	authorizeRespModel := &authorizeRespStruct{}
	json.Unmarshal(body, authorizeRespModel)

	if authorizeRespModel.RedirectUri == "" {
		return "", errors.New("获取失败")
	}
	return authorizeRespModel.RedirectUri, nil
}
