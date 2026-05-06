package scanner

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"net/http"
	"net/url"
	"regexp"

	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/gonmap"
	"github.com/spaolacci/murmur3"
)

type foo2 struct {
	OriginalURL *url.URL
	URL         *url.URL
	response    *gonmap.Response
	req         *http.Request
	client      *http.Client
}

const (
	NotSupportProtocol = "protocol is not support"
)

type URLClient struct {
	*client
	HandlerMatched   func(url *url.URL, final_url *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint, statusCode int)
	HandlerError     func(url *url.URL, err error)
	noRedirectClient *http.Client // 新增：全局禁止跳转的 client
	statusCode       int
}

// 这里是原始的 NewURLScanner 函数，做了修改以支持 JS 跳转处理
// func NewURLScanner(config *Config) *URLClient {
// 	var client = &URLClient{
// 		client:         newConfig(config, config.Threads),
// 		HandlerMatched: func(url *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {},
// 		HandlerError:   func(url *url.URL, err error) {},
// 	}
// 	// 关键：创建一个全局不跟随跳转的 client
// 	// client.noRedirectClient = &http.Client{
// 	// 	Timeout: 10 * time.Second,
// 	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 	// 		return http.ErrUseLastResponse
// 	// 	},
// 	// }
// 	//======= 上面是不跳转逻辑
// 	client.pool.Interval = config.Interval
// 	client.pool.Function = func(in interface{}) {
// 		value := in.(foo2)
// 		URL := value.URL
// 		response := value.response
// 		req := value.req
// 		cli := value.client
// 		if appfinger.SupportCheck(URL.Scheme) == false {
// 			client.HandlerError(URL, errors.New(NotSupportProtocol))
// 			return
// 		}
// 		var banner *appfinger.Banner
// 		var finger *appfinger.FingerPrint
// 		var err error
// 		if response == nil || req != nil || cli != nil {
// 			banner, err = appfinger.GetBannerWithURL(URL, req, cli)

// 			if err != nil {
// 				client.HandlerError(URL, err)
// 				return
// 			}
// 			finger = appfinger.Search(URL, banner)
// 		} else {
// 			//banner, err = appfinger.GetBannerWithResponse(URL, response.Raw, req, cli)
// 			banner, err = appfinger.GetBannerWithURL(URL, req, cli)
// 			if err != nil {
// 				client.HandlerError(URL, err)
// 				return
// 			}
// 			finger = appfinger.Search(URL, banner)
// 			appendTcpBannerInFinger(finger, response.FingerPrint)
// 		}
// 		client.HandlerMatched(URL, banner, finger)
// 	}
// 	return client
// }

// / ============================================================
// 将 http.Header 转为字符串（httpfinger.Banner.Header 通常是 string）
func headerToString(h http.Header) string {
	var buf bytes.Buffer
	for k, vv := range h {
		for _, v := range vv {
			fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
		}
	}
	return buf.String()
}

var titleRegex = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

// 用于提取 favicon link 路径
var faviconRegex = regexp.MustCompile(`(?i)<link[^>]*rel\s*=\s*["']?(?:shortcut )?icon["']?[^>]*href\s*=\s*["']?([^"'\s>]+)["']?`)

// 自定义获取 banner：返回 appfinger.Banner + statusCode + rawBody string
func getBannerWithStatus(URL *url.URL, req *http.Request, cli *http.Client) (*appfinger.Banner, int, string, error) {
	if cli == nil {
		cli = http.DefaultClient
	}

	var httpReq *http.Request
	var err error
	if req != nil {
		httpReq = req
		if httpReq.URL == nil {
			httpReq.URL = URL
		}
	} else {
		httpReq, err = http.NewRequest("GET", URL.String(), nil)
		if err != nil {
			return nil, 0, "", err
		}
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0")
		httpReq.Header.Set("Accept", "*/*")
		httpReq.Header.Set("Connection", "close")
	}

	resp, err := cli.Do(httpReq)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, "", err
	}
	bodyStr := string(bodyBytes)
	headerStr := headerToString(resp.Header)
	// 提取 Title
	title := ""
	if match := titleRegex.FindStringSubmatch(bodyStr); len(match) > 1 {
		title = strings.TrimSpace(htmlUnescape(match[1])) // 可选：处理 &nbsp; 等实体
	}
	// 提取并计算 IconHash（尝试获取 favicon 并计算 mmh3 或 fallback md5）
	iconHash := ""
	faviconURL := findFaviconURL(bodyStr, URL)

	if faviconURL != "" {
		iconHash = fetchAndHashFavicon(faviconURL, cli)
	}
	// 手动构造 appfinger.Banner（不依赖 convBanner）
	banner := &appfinger.Banner{
		Body:   bodyStr,   // 必须是 string
		Header: headerStr, // 大多数指纹规则会用到
		Title:  title,
		Icon:   iconHash,
	}

	return banner, resp.StatusCode, bodyStr, nil
}

// 辅助函数：从 HTML 中提取 favicon 路径
func findFaviconURL(bodyStr string, baseURL *url.URL) string {
	if match := faviconRegex.FindStringSubmatch(bodyStr); len(match) > 1 {
		href := strings.TrimSpace(match[1])
		if u, err := url.Parse(href); err == nil {
			return baseURL.ResolveReference(u).String()
		}
	}
	// fallback: 尝试 /favicon.ico
	return baseURL.Scheme + "://" + baseURL.Host + "/favicon.ico"
}
func fetchAndHashFavicon(faviconURL string, cli *http.Client) string {
	if cli == nil {
		cli = http.DefaultClient
	}

	req, err := http.NewRequest("GET", faviconURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0")

	resp, err := cli.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	iconData, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(iconData) == 0 {
		return ""
	}

	// 1. Base64 编码
	b64 := base64.StdEncoding.EncodeToString(iconData)

	// 2. 每 76 个字符插入 \n，最后再加一个 \n（FOFA 标准）
	var buffer []byte
	for i, ch := range []byte(b64) {
		buffer = append(buffer, ch)
		if (i+1)%76 == 0 {
			buffer = append(buffer, '\n')
		}
	}
	buffer = append(buffer, '\n')

	// 3. 对带换行的 Base64 字节流计算 MurmurHash3 32位
	hash32 := murmur3.Sum32(buffer)

	// 4. 转为有符号 int32 并返回字符串形式（如 "-1507567067"）
	iconHash := int32(hash32)
	return fmt.Sprintf("%d", iconHash)
}

// 可选：简单 HTML 实体转义处理
func htmlUnescape(s string) string {
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return s
}
func NewURLScanner(config *Config) *URLClient {
	sc := &URLClient{
		client: newConfig(config, config.Threads),
		HandlerMatched: func(url *url.URL, final_url *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint, statusCode int) {
		},
		HandlerError: func(url *url.URL, err error) {},
	}

	sc.pool.Interval = config.Interval
	sc.pool.Function = func(in interface{}) {
		value := in.(foo2)
		OriginalURL := value.OriginalURL
		origURL := value.URL
		//response := value.response
		req := value.req
		cli := value.client
		// 强制使用不跳转的 client
		//cli := sc.noRedirectClient
		if value.client != nil {
			cli = value.client
		}

		if !appfinger.SupportCheck(origURL.Scheme) {
			//sc.HandlerError(origURL, errors.New(NotSupportProtocol))
			return
		}

		// 第一次请求
		banner, statusCode, bodyStr, err := getBannerWithStatus(origURL, req, cli)
		if err != nil {
			sc.HandlerError(origURL, err)
			return
		}
		// fmt.Println("第一次请==============================")
		// fmt.Printf("%d\ttitle:%s\n", statusCode, banner.Title)
		// fmt.Printf("%s\n", banner.Header)
		// fmt.Printf("%s\n", banner.Body)
		// fmt.Printf("%s\n", banner.Icon)
		// fmt.Println("==============================")
		finalURL := origURL
		finalBanner := banner
		//response := value.response
		//指纹识别
		finger := appfinger.Search(finalURL, finalBanner)
		//合并 TCP 指纹
		// if response != nil {
		// 	fmt.Printf("finger %s, response.FingerPrint %s\n", finger, response.FingerPrint)
		// 	appendTcpBannerInFinger(finger, response.FingerPrint)
		// }

		//回调拿到指纹
		sc.HandlerMatched(OriginalURL, finalURL, finalBanner, finger, statusCode)
		//fmt.Println("第一次 但是这里也不知道第一次请求后是否存在指纹吗，因为是回调函数")
		redirectPath := ""

		var jsRedirectRegex = regexp.MustCompile(`(?i)\w+\.location(?:\.href|\.replace|\.assign)?\s*(?:\(|=)\s*["']([^"']+?)["']\s*(?:[);]|$)`)

		var metaRefreshRegex = regexp.MustCompile(`(?i)<meta[^>]+http-equiv\s*=\s*["']?refresh["']?[^>]+content\s*=\s*["']?\d+\s*;\s*url\s*=\s*['"]?([^'"\s>]+)`)
		var match []string
		if statusCode == 200 {
			if match = jsRedirectRegex.FindStringSubmatch(bodyStr); len(match) > 1 {
				redirectPath = match[1]
				// 检查是否以 ../ 开头
				if strings.HasPrefix(redirectPath, "../") {
					//fmt.Println("这里 ../ 跳过了", redirectPath)
					return // 或者 continue
				}
				//fmt.Printf("[jsRedirectRegex 跳转] 检测到: %s\n", redirectPath) // 可选调试日志
			} else if match := metaRefreshRegex.FindStringSubmatch(bodyStr); len(match) > 1 {
				redirectPath = match[1]
				//fmt.Printf("[meta refresh 跳转] 检测到: %s\n", redirectPath) // 可选调试日志
			}
		}
		if redirectPath == "" || redirectPath == "/html/ie.html" ||
			redirectPath == "resourcec/html/ie.html" {
			//fmt.Println("redirectPath 为 nil 跳过")
			return
		}
		if redirectPath != "" {
			if strings.HasPrefix(strings.ToLower(redirectPath), "http://") ||
				strings.HasPrefix(strings.ToLower(redirectPath), "https://") ||
				strings.HasPrefix(redirectPath, "//") {
				//fmt.Println("跳过绝对路径跳转:", redirectPath)
			} else if redirectURL, parseErr := url.Parse(redirectPath); parseErr == nil {
				targetURL := origURL.ResolveReference(redirectURL)

				// 第二次请求目标页面
				//fmt.Println("第二次请求目标页面")
				if targetBanner, _, _, targetErr := getBannerWithStatus(targetURL, nil, sc.noRedirectClient); targetErr == nil {
					finalBanner = targetBanner
					finalURL = targetURL

					fakeFinger := new(appfinger.FingerPrint)
					sc.HandlerMatched(OriginalURL, finalURL, finalBanner, fakeFinger, statusCode)
					//fmt.Printf("第二次请求目标页面 %-60s\n", finalURL)
					// fmt.Printf("第二次请求目标页面 body %s\n", banner.Body)
				}
				// 失败时回退到原页面
			}
		}

		//finger := appfinger.Search(finalURL, finalBanner)
		// 合并 TCP 指纹
		// if response != nil {
		// 	fmt.Printf("finger %s, response.FingerPrint %s\n", finger, response.FingerPrint)
		// 	appendTcpBannerInFinger(finger, response.FingerPrint)
		// }

		// 这里是回调 然后找指纹特征
		//fmt.Printf("这是js匹配后的扫描结果 \n")

	}
	return sc
}

func (c *URLClient) Push(URL *url.URL, response *gonmap.Response, req *http.Request, client *http.Client) {
	c.pool.Push(foo2{
		OriginalURL: URL,
		URL:         URL, // 初始时相同
		response:    response,
		req:         req,
		client:      client,
	})
}

func appendTcpBannerInFinger(finger *appfinger.FingerPrint, gonmapFinger *gonmap.FingerPrint) *appfinger.FingerPrint {
	if gonmapFinger.ProductName != "" {
		if gonmapFinger.Version != "" {
			finger.AddProduct(gonmapFinger.ProductName + "/" + gonmapFinger.Version)
		}
		finger.AddProduct(gonmapFinger.ProductName)
	}

	if gonmapFinger.OperatingSystem != "" {
		finger.AddProduct(gonmapFinger.OperatingSystem)
	}

	if gonmapFinger.DeviceType != "" {
		finger.AddProduct(gonmapFinger.DeviceType)
	}

	if gonmapFinger.Info != "" {
		finger.AddProduct(gonmapFinger.Info)
	}

	if gonmapFinger.Hostname != "" {
		finger.Hostname = gonmapFinger.Hostname
	}
	return finger
}
