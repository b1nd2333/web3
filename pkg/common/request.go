package common

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/andybalholm/brotli"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// NewHTTPClientWithProxy 根据代理创建client
func NewHTTPClientWithProxy(proxyAddress string) (*http.Client, error) {
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if proxyAddress != "" {
		// 解析代理地址
		proxyURL, err := url.Parse(proxyAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy address: %v", err)
		}

		if proxyURL.Scheme == "socks5" {
			// 设置 SOCKS5 代理并进行身份验证
			dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
			if err != nil {
				return nil, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
			}
			// 创建 HTTP Transport 使用 SOCKS5 代理
			transport = &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 忽略 HTTPS 错误
				},
			}
		} else {
			transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 忽略 HTTPS 错误
				},
			}
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return client, nil
}

// DecompressBody 解压请求包
func DecompressBody(body []byte) []byte {
	// 先尝试 Brotli 解压
	brReader := brotli.NewReader(bytes.NewReader(body))
	if decompressed, err := io.ReadAll(brReader); err == nil {
		return decompressed
	}

	// 再尝试 Gzip 解压
	gzReader, err := gzip.NewReader(bytes.NewReader(body))
	if err == nil {
		defer gzReader.Close()
		if decompressed, err := io.ReadAll(gzReader); err == nil {
			return decompressed
		}
	}

	return body
}
