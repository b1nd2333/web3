package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b1nd2333/web3/pkg/common"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BindDCRespStruct struct {
	Location string `json:"location"`
}

func BindDCLink(proxyStr string, dcToken string, bindUri string) (string, error) {
	// 创建 HTTP 客户端
	client, err := common.NewHTTPClientWithProxy(proxyStr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", bindUri, strings.NewReader("{\"permissions\":\"0\",\"authorize\":true,\"integration_type\":0,\"location_context\":{\"guild_id\":\"10000\",\"channel_id\":\"10000\",\"channel_type\":10000},\"dm_settings\":{\"allow_mobile_push\":false}}"))
	if err != nil {
		return "", err
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}
	req.Header.Set("User-Agent", userAgents[time.Now().UnixNano()%int64(len(userAgents))])
	req.Header.Add("Authorization", dcToken)
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("x-super-properties", "eyJvcyI6Ik1hYyBPUyBYIiwiYnJvd3NlciI6IkNocm9tZSIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbiIsImhhc19jbGllbnRfbW9kcyI6ZmFsc2UsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMzAuMC4wLjAgU2FmYXJpLzUzNy4zNiIsImJyb3dzZXJfdmVyc2lvbiI6IjEzMC4wLjAuMCIsIm9zX3ZlcnNpb24iOiIxMC4xNS43IiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfY3VycmVudCI6IiIsInJlZmVycmluZ19kb21haW5fY3VycmVudCI6IiIsInJlbGVhc2VfY2hhbm5lbCI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjQxMDcwNiwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbCwiY2xpZW50X2xhdW5jaF9pZCI6IjMxMmJhZDI0LWU3ZGEtNDQyOS04YTMxLTA4NTZhNzliYTFiZCIsImNsaWVudF9hcHBfc3RhdGUiOiJ1bmZvY3VzZWQifQ==")
	req.Header.Add("x-debug-options", "bugReporterEnabled")
	req.Header.Add("x-discord-locale", "en-US")
	req.Header.Add("content-type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return "", err
	}

	bindDCRespModel := &BindDCRespStruct{}
	decompressBody := common.DecompressBody(body)
	_ = json.Unmarshal(decompressBody, bindDCRespModel)
	return bindDCRespModel.Location, err
}
