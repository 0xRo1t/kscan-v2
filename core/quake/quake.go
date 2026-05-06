package quake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kscan/core/slog"
	"kscan/lib/color"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

type QuakeRequest struct {
	Query       string   `json:"query"`
	Start       int      `json:"start"`
	Size        int      `json:"size"`
	IgnoreCache bool     `json:"ignore_cache"`
	Include     []string `json:"include,omitempty"`
}

type QuakeData struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Hostname   string `json:"hostname"`
	Time       string `json:"time"`
	Components []struct {
		ProductNameCn string `json:"product_name_cn"`
		ProductNameEn string `json:"product_name_en"`
		Version       string `json:"version"`
	} `json:"components,omitempty"`
	Service struct {
		Name     string `json:"name"`
		Cert     string `json:"cert,omitempty"`
		Response string `json:"response,omitempty"`
		HTTP     struct {
			Title           string `json:"title"`
			Host            string `json:"host"`
			Body            string `json:"body,omitempty"`
			RequestHeaders  string `json:"request_headers,omitempty"`
			ResponseHeaders string `json:"response_headers,omitempty"`
			Icp             struct {
				Licence     string `json:"licence"`
				UpdateTime  string `json:"update_time"`
				IsExpired   bool   `json:"is_expired"`
				LeaderName  string `json:"leader_name"`
				Domain      string `json:"domain"`
				MainLicence struct {
					Licence string `json:"licence"`
					Unit    string `json:"unit"`
					Nature  string `json:"nature"`
				} `json:"main_licence"`
				ContentTypeName string `json:"content_type_name"`
				LimitAccess     bool   `json:"limit_access"`
			} `json:"icp,omitempty"`
		} `json:"http"`
		TLS struct {
			ServerCertificates []struct {
				Certificate struct {
					Parsed struct {
						Subject struct {
							CommonName []string `json:"common_name"`
						} `json:"subject"`
					} `json:"parsed"`
				} `json:"certificate"`
			} `json:"server_certificates,omitempty"`
		} `json:"tls,omitempty"`
	} `json:"service"`
	Location struct {
		CountryEn  string `json:"country_en"`
		Province   string `json:"province_cn"`
		City       string `json:"city_cn"`
		DistrictCn string `json:"district_cn"`
		DistrictEn string `json:"district_en"`
		Owner      string `json:"owner"`
	} `json:"location"`
}

type QuakeResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []QuakeData `json:"data"`
	Meta    struct {
		Pagination struct {
			Total int `json:"total"`
			Count int `json:"count"`
		} `json:"pagination"`
	} `json:"meta"`
}

type QuakeClient struct {
	ApiKey string
	URL    string
}

var quakeClient *QuakeClient
var quakeResults []QuakeData

func NewQuakeClient(apiKey string) *QuakeClient {
	return &QuakeClient{
		ApiKey: apiKey,
		URL:    "https://quake.360.net/api/v3/search/quake_service",
	}
}

func (c *QuakeClient) Search(query string, size int) (*QuakeResponse, error) {
	payload := QuakeRequest{Query: query, Start: 0, Size: size, IgnoreCache: false}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-QuakeToken", c.ApiKey)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}
	var quakeResp QuakeResponse
	json.Unmarshal(body, &quakeResp)
	return &quakeResp, nil
}

func Init(apiKey string) {
	quakeClient = NewQuakeClient(apiKey)
}

func Run(query string, size int) {
	result, err := quakeClient.Search(query, size)
	if err != nil {
		slog.Println(slog.DATA, fmt.Sprintf("[quake] 查询失败: %v", err))
		os.Exit(1)
	}
	slog.Println(slog.DATA, fmt.Sprintf("[quake] 本次搜索关键字为：%v", color.Red(query)))
	quakeResults = result.Data
	displayResponse(quakeResults, result.Meta.Pagination.Total)
}

func displayResponse(data []QuakeData, total int) {
	seen := make(map[string]bool)
	count := 0
	for _, item := range data {
		key := fmt.Sprintf("%s:%d", item.IP, item.Port)
		if seen[key] {
			continue
		}
		seen[key] = true
		count++
		title := item.Service.HTTP.Title
		icp := item.Service.HTTP.Icp.Licence
		domain := item.Service.HTTP.Host
		//product := ""
		if len(item.Components) > 0 {
			//product = item.Components[0].ProductNameEn
		}
		city := item.Location.City
		district := item.Location.DistrictCn
		servers := fixServiceName(item.Service.Name)
		host := fmt.Sprintf("%s://%s:%d", servers, item.IP, item.Port)
		slog.Println(slog.DATA, fmt.Sprintf("%s | %s | %s %s | %s | %s",
			padRight(host, 30),
			padRight(icp, 20),
			padRight(city, 7),
			padRight(district, 7),
			padRight(strings.TrimSpace(title), 25),
			padRight(domain, 25),
		))
	}
	slog.Println(slog.DATA, fmt.Sprintf("[quake] 命中总数: %d, 去重后实际条数: %d", total, count))
}

func fixServiceName(name string) string {
	name = strings.ToLower(name)

	if strings.Contains(name, "(http)") {
		return "http"
	}
	if strings.Contains(name, "(https)") {
		return "https"
	}

	if name == "http/ssl" || strings.Contains(name, "https") {
		return "https"
	}
	if strings.Contains(name, "http") {
		return "http"
	}

	return name
}

func GetUrlTarget() []string {
	var targets []string
	for _, item := range quakeResults {
		servers := fixServiceName(item.Service.Name)
		targets = append(targets, fmt.Sprintf("%s://%s:%d", servers, item.IP, item.Port))
	}
	return uniqueStringSlice(targets)
}

func GetHostTarget() []string {
	var hosts []string
	for _, item := range quakeResults {
		hosts = append(hosts, item.IP)
	}
	return uniqueStringSlice(hosts)
}

func uniqueStringSlice(in []string) []string {
	set := make(map[string]struct{})
	var out []string
	for _, s := range in {
		if _, ok := set[s]; ok {
			continue
		}
		set[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// padRight 按终端显示宽度（中文占2格）右侧补空格到指定宽度
func padRight(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}
