package fofa

import (
	"fmt"
	"kscan/app"
	"kscan/core/slog"
	"kscan/lib/color"
	"kscan/lib/fofa"
	"kscan/lib/misc"
	"regexp"
	"strconv"
	"strings"
)

var this *fofa.Client
var keywordSlice []string
var results []fofa.Result

func Init(email, key string) {
	//设置日志输出器
	fofa.SetLogger(slog.Debug())
	//初始化fofa模块
	this = fofa.New(email, key)
	this.SetSize(app.Setting.FofaSize)
	//获取所有关键字
	keywordSlice = makeKeywordSlice()
}

func Run() {
	//对每个关键字进行查询
	for _, keyword := range keywordSlice {
		//slog.Printf(slog.WARN, "本次搜索关键字为：%v", keyword)
		slog.Println(slog.DATA, fmt.Sprintf("[+] 本次搜索关键字为：%v", color.Red(keyword)))
		size, r := this.Search(keyword)
		displayResponse(r)
		slog.Printf(slog.DATA, "[+] 本次搜索，返回结果总条数为：%d,此次返回条数为：%d", size, len(r))
		results = append(results, r...)
	}
}

func makeKeywordSlice() []string {
	var keywordSlice []string
	if app.Setting.FofaFixKeyword == "" {
		keywordSlice = app.Setting.Fofa
	} else {
		for _, keyword := range app.Setting.Fofa {
			keyword = strings.ReplaceAll(app.Setting.FofaFixKeyword, "{}", keyword)
			keywordSlice = append(keywordSlice, keyword)
		}
	}
	return keywordSlice
}

func GetUrlTarget() []string {
	var strSlice []string
	for _, r := range results {
		Fix(&r)
		strSlice = append(strSlice, r.Host)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

func GetHostTarget() []string {
	var strSlice []string
	for _, r := range results {
		strSlice = append(strSlice, r.Ip)
	}
	strSlice = misc.RemoveDuplicateElement(strSlice)
	return strSlice
}

func displayResponse(r []fofa.Result) {
	for _, row := range r {
		Fix(&row)
		m := row.Map()
		m["Header"] = ""
		m["Cert"] = ""
		m["Title"] = row.Title
		m["Host"] = ""
		m["As_organization"] = ""
		if m["Domain"] == "" {
			m["Ip"] = ""
		}
		m["Port"] = ""
		m["Country_name"] = ""
		m = misc.FixMap(m)
		if m["Banner"] != "" {
			m["Banner"] = misc.FixLine(m["Banner"])
			m["Banner"] = misc.StrRandomCut(m["Banner"], 20)
		}
		//locInfo := strings.Trim(fmt.Sprintf("%s, %s, %s", row.As_organization, row.Country_name, row.City), ", ")
		locInfo := strings.Trim(fmt.Sprintf("%s, %s", row.Country_name, row.City), ", ")
		line := fmt.Sprintf("%-30v %-"+strconv.Itoa(misc.AutoWidth(row.Title, 26))+"v %v",
			row.Host,
			row.Title,
			//color.StrMapRandomColor(m, app.Setting.CloseColor, []string{"Server"}, []string{}),
			locInfo,
		)
		slog.Println(slog.DATA, line)
	}
}

func Fix(r *fofa.Result) {
	//修复title
	if r.Title == "" && r.Protocol != "" {
		r.Title = strings.ToUpper(r.Protocol)
	}
	r.Title = misc.FixLine(r.Title)
	//修改host
	if r.Host == "" {
		r.Host = r.Ip
	}
	// 这里的 r.Host 就是 ip:port 的形式

	if strings.Contains(strings.ToLower(r.Protocol), "nacos(http)") {
		r.Host = "http://" + r.Host + "/nacos"
	} else {
		// 4. 如果 Protocol 不含 http，再走常规的协议头补全逻辑
		if !regexp.MustCompile(`(?i)\w+://`).MatchString(r.Host) {
			protocol := r.Protocol
			if protocol == "" {
				protocol = "http"
			}
			r.Host = protocol + "://" + r.Host
		}
	}

}
