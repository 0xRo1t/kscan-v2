package finger

import (
	"embed"
	"encoding/json"
	"log"
	"os"
)

// Fingerprint 结构体定义了单个指纹的规则
type Fingerprint struct {
	Cms      string   `json:"cms"`
	Keyword  []string `json:"keyword"`
	Location string   `json:"location"`
	Method   string   `json:"method"`
}

// 结构体用于解析整个指纹文件
type FingerprintData struct {
	Fingerprints []Fingerprint `json:"fingerprint"`
}

var Fingerprints []Fingerprint

// 声明一个私有的 embed.FS 变量，用于存储嵌入的文件系统
var embeddedFS embed.FS

// SetEmbeddedFS 供外部包（如 main/InitKscan）调用，用于设置嵌入的文件系统
func SetEmbeddedFS(fs embed.FS) {
	embeddedFS = fs
}

// ... 你的 FingerprintData, Fingerprints 等定义 ...

func LoadFingerJSON(path string) {
	var file []byte
	var err error

	// **重点修改区域：**
	// 1. 优先尝试从嵌入的文件系统（embeddedFS）中读取。
	// 2. path 参数在这里被用作嵌入式文件系统中的文件名。
	if embeddedFS != (embed.FS{}) {
		file, err = embeddedFS.ReadFile(path)
		if err == nil {
			// 如果从嵌入FS中读取成功，则跳过os.ReadFile
			goto process
		}
		// 如果嵌入式读取失败，可能是因为文件不存在或路径错误，继续尝试读取磁盘文件。
	}

	// 回退到原始的磁盘文件读取逻辑
	file, err = os.ReadFile(path)
	if err != nil {
		log.Fatalf("无法读取指纹文件 [%s]： %v", path, err)
	}

process: // 跳转标签，用于统一处理读取到的字节数据
	var fpData FingerprintData
	if err := json.Unmarshal(file, &fpData); err != nil {
		log.Fatalf("无法解析指纹文件: %v", err)
	}
	Fingerprints = fpData.Fingerprints
}
