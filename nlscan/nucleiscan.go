package nlscan

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"kscan/core/slog"
	"kscan/lib/color"
	"path/filepath"
	"sync"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

//go:embed poc/**
var pocEmbedFS embed.FS

// 全局变量
var once sync.Once // 保证初始化只执行一次

func Init_poc(ne *nuclei.NucleiEngine) {
	// 1. 在外部定义计数器
	var successCount int
	var totalFiles int
	// 第一步：收集所有嵌入的 POC 模板内容（作为 []string 或 []byte）

	// 直接从嵌入的文件系统加载模板：遍历 poc 下所有文件并通过 ne.ParseTemplate 加载
	_ = fs.WalkDir(pocEmbedFS, "poc", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// 只处理 yaml/yml 文件
		ext := filepath.Ext(path)
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}
		totalFiles++ // 记录扫描到的符合后缀的文件总数
		b, err := pocEmbedFS.ReadFile(path)
		if err != nil {
			return err
		}
		// ParseTemplate 会将模板解析并注册到引擎的 parser/cache 中
		if _, err := ne.ParseTemplate(b); err != nil {
			// `打印错误但继续尝试加载其他模板`
			//fmt.Printf("解析模板失败 %s: %v\n", path, err)
		} else {
			// 2. 解析成功则计数器加一
			successCount++
		}
		return nil
	})
	fmt.Printf("[+] 共发现 %d 个POC文件, 成功加载 %d 个模板\n", totalFiles, successCount)
}

func Nucleiscan(target []string, tags []string) {

	//startTime = time.Now()
	ctx := context.Background()
	// 创建引擎（不通过磁盘路径加载模板）
	ne, err := nuclei.NewNucleiEngineCtx(ctx,
		// 这里是为了指向不存在的目录，然后只加载我们自己的poc
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{Templates: []string{""}}),
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Severity: "critical,high,medium",
			//Tags:     strings.Split("solr", ",")
			Tags:        tags,
			ExcludeTags: []string{}, // 设置为空，覆盖默认的ignore文件规则
		}),
		//nuclei.WithGlobalRateLimitCtx(ctx, 200, time.Second), // 全局限速（可选）
		nuclei.DisableUpdateCheck(),   // 禁用更新检查（可选）
		nuclei.DisableTemplateCache(), // 禁用模板缓存（可选）
	)
	once.Do(func() {
		Init_poc(ne) // 使用闭包捕获 ne
	})
	nuclei.DefaultConfig.GetIgnoreFilePath()
	if err != nil {
		panic("引擎初始化失败: " + err.Error())
	}
	defer ne.Close()
	// 3. 添加目标
	// targets := []string{
	// 	"http://123.58.224.8:33328/",
	// }
	ne.LoadTargets(target, false)
	ne.Options().MatcherStatus = true
	// 4. 开始扫描 + 美观输出
	err = ne.ExecuteCallbackWithCtx(ctx, func(event *output.ResultEvent) {

		if event == nil {
			return
		}
		if event.Matched == "" {
			// fmt.Printf("[-] [%s] [%s] [%s] [%s]\n",
			// 	event.Info.SeverityHolder.Severity.String(),
			// 	event.TemplateID,
			// 	event.Info.Tags,
			// 	event.Host)
		} else {
			// fmt.Printf("[%s] [%s] %s\n",
			// 	color.Red(event.Info.SeverityHolder.Severity.String()),
			// 	color.Green(event.TemplateID),
			// 	color.RedB(event.Matched),
			// )

			//msg := fmt.Sprintf("[vuln] [%s] [%s] %s", event.Info.SeverityHolder.Severity.String(), event.TemplateID, event.Matched)

			msg1 := fmt.Sprintf("[vuln] [%s] [%s] %s",
				color.Red(event.Info.SeverityHolder.Severity.String()),
				color.Green(event.TemplateID),
				color.RedB(event.Matched),
			)
			// 控制台（带颜色）
			// fmt.Printf("[%s] [%s] %s\n",
			// 	color.Red(event.Info.SeverityHolder.Severity.String()),
			// 	color.Green(event.TemplateID),
			// 	color.RedB(event.Matched),
			// )

			slog.Println(slog.DATA, msg1)
			slog.Println(slog.DATA, "================  请求数据包 ================")
			slog.Println(slog.DATA, event.Request)
			slog.Println(slog.DATA, "================  请求响应包 ================")
			slog.Println(slog.DATA, event.Response)
			slog.Println(slog.DATA, "============================================")
		}
	})
	if err != nil {
		//fmt.Printf("⚠️ 警告：扫描执行失败，已跳过。错误信息: %s\n", err.Error())
		//panic("扫描执行失败: " + err.Error())
		return
	}
}
