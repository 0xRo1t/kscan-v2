fmt.Printf("hearder %s\n", headers)

## 设置的 不跳转逻辑在
func NewURLScanner

## 匹配指纹
func match_finger

## 漏扫的地方
func generateURLScanner

# 调用web扫描
URLScanner = generateURLScanner(wg)

## 定义参数
func (o *args) define()

## 输入目标后 -t。自动扫描端口的逻辑
if app.Setting.Check == false {
			pushTarget(netloc)
		}

## 加了这一段防止重复扫描  run.go 里面的逻辑
urlKey := URL.String()
	if _, loaded := urlFilter.LoadOrStore(urlKey, true); loaded {
		// 如果已经存在，直接跳过，不再向下分发
		return
	}

## 输入 -t 默认扫端口的行为
netloc, port := uri.SplitWithNetlocPort(expr)
		if uri.IsIPv4(netloc) {
			PortScanner.Push(net.ParseIP(netloc), port)
		}

## 扫描器流程
输入  ./kscan -t http://180.101.183.177:8089
然后进入  if uri.IsURL(expr)
之后 if uri.IsIPv4(expr) 在这里扫端口了


所以不想扫端口，直接写成这样即可
if uri.IsIPv4(expr) {
		// 这是默认扫描端口的行为
		fmt.Printf("走了这里扫端口的000\n")
		//IPScanner.Push(net.ParseIP(expr))
		// if app.Setting.Check == true {
		// 	pushURLTarget(uri.URLParse("http://"+expr), nil)
		// 	pushURLTarget(uri.URLParse("https://"+expr), nil)
		// }
		return
	}
## 但是如果这样写，./kscan -t 180.101.183.177  这样就不会输入了,加了端口也不会输出了
./kscan -t 180.101.183.177 -p 8080,8089,8090 这样会直接退出程序，什么都没做

扫描两次的逻辑在 
func NewURLScanner(config *Config) *URLClient 这里因为是存在 js跳转 所以扫了两次，需要优化下

## func IsNetlocPort 这个函数是检查
是否是 [Domain or IP]:Port 这种形式的

## url的扫描逻辑
./kscan -t http://180.101.183.177:8089
直接走 uri.IsURL(expr) -> 然后指纹扫描/js路径扫描