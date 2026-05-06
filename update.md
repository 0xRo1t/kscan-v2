## 20260325 更新了fofa查询结果的格式
locInfo := strings.Trim(fmt.Sprintf("%s, %s, %s", row.Country, row.Country_name, row.City), ", ")
http://sh.hbgs.com.cn          河北高速公路集团有限公司石黄分公司 CN, 中国, Xuzhou


## 20260318更新 默认输出文件
输出的时候这样写就好了，然后输出文件名字为当前时间
slog.Println(slog.DATA, fmt.Sprintf("%-40s 🤓 %s", final_url, r.Matches))

## 全局忽略证书
func Init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

## 删除报错
protocol is not support
注释掉：
//sc.HandlerError(origURL, errors.New(NotSupportProtocol))

报错：
fmt.Printf("[ERROR] %s: %v\n", url, err)
##
http://202.100.163.17:5050                                        |🤓|
这里 ../ 跳过了 ../WebReport/ReportServer?op=fs
http://58.210.216.117                                             |🤓| finereport
http://58.210.216.117/FineReport/ReportServer?op=fs               |🤓| finereport
## 修复 //go:embed .nuclei-ignore
内嵌了 .nuclei-ignore

[bug]
http://112.15.120.59:8888/seeyon/autoinstall/zxsetup.dmg |👀| 
这里是一个安装包，但是没扫出来指纹？？

[bug] 
http://101.69.122.168:3120 |👀| 有WAF
[*]2025/12/25 09:53:41 无法识别的Target字符串:ncacn_http://101.67.62.95:8009

其中 无法识别的Target字符串 这句话是打印出来的， 这里很奇怪 ncacn不知道是什么东西

## [已修复] 还是重复扫描  if match == nil return  添加了这个逻辑，没匹配到就不扫描了 减少了扫描的次数
./kscan -t http://180.101.183.177:8089
[+] 加载 20587 条指纹
[+]2025/12/22 17:57:49 Domain、IP、Port、URL、Hydra引擎已准备就绪
http://180.101.183.177:8089 |👀| e-office;Bootstrap;fanwei;php 
http://180.101.183.177:8089 |👀| e-office;Bootstrap;fanwei;php 

## [bug] 还是跳转的问题
http://180.101.183.177:8089/
这个链接在body 302跳转了  var redirectUrl = "/general/login/index.php" ;

之后跳转的这个链接 存在一个跳转
<form id='loginform' action='/general/login/index.php?c=Login&a=loginCheck' 
就又跳转到这里了，好奇怪啊


60.175.49.104:8010 这个站，跳转之后，是存在指纹的，但是 title是空的

更新了，statusCode != 200 的时候才强制跳转  但是发现这样出问题了，应该改的思路是 ，第一次是200 

http://139.9.203.249:9088    这个会跳转到 /true 更奇怪了
## [已修复] 跳转出现bug了   判断了js匹配如果有 http前缀就不跳转了
现在是默认跳转，然后跳转之后如果是200，就还匹配这些跳转链接，导致了这个问题

可以修嘎一下逻辑，先去扫描状态吗是200的 这个不跳转，然后再去扫描302的 这个跟随跳转就好了

http://astwy.cn:8989/ 这个站，正常访问302 跳转到
Location: /login.aspx?ReturnUrl=%2fdefault.aspx
之后就是 landray-eis的指纹了， 但是他又跳转到 逆天

window.location.href="https://oapi.dingtalk.com/connect/oauth2/sns_authorize?appid="


================
[high] [Apusic-lfi] http://118.190.134.163:8080/admin/protected/selector/server_file/files?folder=/

## 自实现了 fofa的 ico hash函数
fetchAndHashFavicon
## 
[meta refresh 跳转] 检测到: index.jsp
http://222.134.77.42:8089/index.jsp                         yonyou nc 
[meta refresh 跳转] 检测到: index.jsp
[meta refresh 跳转] 检测到: index.jsp
http://114.242.114.169:8087/login.html                      ufida 
http://139.9.234.31:8090/index.jsp                          yonyou nc 
http://106.124.147.176:9006/index.jsp                       ufida 
[meta refresh 跳转] 检测到: index.jsp
http://220.194.141.190:9999/index.jsp                       yonyou nc 
[meta refresh 跳转] 检测到: index.jsp
http://222.184.237.179:9999/index.jsp                       yonyou nc  这个 应该是 匹配到ufida才对
[meta refresh 跳转] 检测到: index.jsp
http://219.138.90.213:8089/index.jsp                        yonyou nc 

#  139.196.183.21:8009 没有匹配到跳转到
状态吗是200 但是 跳转到js是 ;url=index.jsp
发现了一个特征，有一些跳转没有用 location 用的是
<meta http-equiv=refresh content=0;url=index.jsp>
或者
<META HTTP-EQUIV="refresh" CONTENT="0;URL=index.jsp">
<meta http-equiv="refresh" content="0;url=index.jsp">

## 改了一下输出的逻辑 fmt.Printf("%-60s", finalURL) 然后拼接 fmt.Printf("%s \n", r.Matches)
这样输出指纹，就是最终的url了， 而不是基础的url

## 2025.12.15 优化了 状态吗是200 但是会跳转到情况
## 这下之后 指纹都直接用那个 跳转后到就好了
但是现在还有个问题，就是会重复扫描，这里不知道是为什么重复扫描了
./kscan -t http://121.46.30.165:8085 -p 8085
[+] 加载 20580 条指纹
http://121.46.30.165:8085                                   Struts2;seeyon-DEE可视化配置工具 
[high] [seeyon_dee_weakpasswd_dee_admin] http://121.46.30.165:8085/dee/login!checkLogin.do
http://121.46.30.165:8085                                   Struts2;seeyon-DEE可视化配置工具 
[high] [seeyon_dee_weakpasswd_dee_admin] http://121.46.30.165:8085/dee/login!checkLogin.do

http://180.159.27.31:7070                                    
http://222.178.41.222:6688                                   
http://47.101.178.187:18086                                  
http://wshine.net.cn:8071                                    
http://116.198.46.169:993                                    
http://erp.kgk.com.cn:8072                                   
http://43.130.249.241:1234 
一个都没有扫到 u8-crm
 HTTP/1.1 200 OK
Server: Apache/2.4.41 (Win32) PHP/5.4.38
X-Powered-By: PHP/5.4.38
Content-Length: 233
Content-Type: text/html; charset=UTF-8
Date: Mon, 15 Dec 2025 03:30:32 GMT

响应包是 200 导致了没有跟随跳转
===========================
<script type="text/javascript">
					var winPath = "top";
					if (eval("window." + winPath))
						eval("window." + winPath).location="/login/login.php";
					else
						window.top.location="/login/login.php";
				</script>                      
# 
添加了 accept头 ，如果不添加比如 56.155.135.199:9090 就会显示401 导致无法访问指纹
req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")


http://3.29.62.226:8085                                     QNAP-NAS 
https://40.177.166.36:58000                                 SAPNetWeaver 
http://64.119.205.118:8085                                  Struts2;Tomcat 
http://212.174.69.242:5000                                  Struts2;Tomcat 
http://181.115.48.42:8085                                   Struts2;Tomcat 
http://63.143.90.18:8085                                    jboss;RedHat JBoss;Struts2;Tomcat 
http://63.143.89.190:8085                                   Struts2;Tomcat 
https://3.252.26.188:2850                                   SAPNetWeaver 
这几个站 是 jboss但是没扫到
发现原因是 func match_finger 函数的开头 header 里面是原始的 header 然后有大小写的区分
body的话目前是转小写了， 所以body尽量写小写的比较好
## 设置的 不跳转逻辑在
func NewURLScanner
如果想修改 跳转，直接注释掉那一段即可，设置了注释的

## 当前扫描器行为
默认不跟随跳转， 不跟随跳转 很多细分的指纹无法匹配到，比如说
http://60.21.224.26:5000/ 这个 访问后是 302 站跳转后是 致远A6+协同管理软件 V7.0SP1
但因为没有做跳转，我们抓到的特征只是 Location: /seeyon/index.jsp 导致了无法判断具体的版本
正常来说 应该匹配到指纹 seeyon-A6+
但是莫名其妙的 扫到了其他的漏洞哈哈哈哈 其他的 seeyon tags的漏洞还是好抽象

fmt.Printf("hearder %s\n", headers)

## 这个站还是因为 header匹配不到 loaction 导致无法匹配到指纹
113.204.227.162:8885

## 误报问题
[high] [seeyon_dee_weakpasswd_dee_admin] http://121.11.101.158:8085/dee/login!checkLogin.do
[high] [seeyon_dee_weakpasswd] http://121.11.101.158:8085/dee/login!checkLogin.do
http://113.125.69.44:8085                                   seeyon-DEE可视化配置工具 
http://123.60.136.11:8085                                   seeyon-DEE可视化配置工具 
http://121.36.70.187:8085                                   seeyon-DEE可视化配置工具 
http://58.137.205.123:8085                                  seeyon-DEE可视化配置工具 
[high] [seeyon_dee_weakpasswd] http://47.104.132.78:8888/dee/login!checkLogin.do
[high] [seeyon_dee_weakpasswd_dee_admin] http://47.104.132.78:8888/dee/login!checkLogin.do
真的 但是扫了两遍 [high] [seeyon_dee_weakpasswd_dee_admin] http://121.46.30.165:8085/dee/login!checkLogin.do
误报 [high] [seeyon_dee_weakpasswd] http://121.46.30.165:8085/dee/login!checkLogin.do


## [已修复]关于 抓取 body和 header的时候

有时候bp抓到的 和 我们的代码抓到的 不一样，是因为cookie不一致导致的
还有就是 访问同一个站点 这里举例子： http://121.36.70.187:8085
可以看到 代码抓不到 location头 还有bp抓到的默认都是带cookie访问的结果，一般不会有Set-Cookie字段，这里是手动删除了cookie的结果
================ 下面是bp抓到的 ====================
HTTP/1.1 302 
Set-Cookie: deeClientSeesionId=271DE82F4BC9111D02CC9509096EF1A6; Path=/; HttpOnly
Location: /dee/login!welcome.do
Content-Type: text/html;charset=UTF-8
Content-Length: 0
Date: Thu, 11 Dec 2025 04:09:31 GMT
================ 代码抓到的 ====================
HTTP/1.1 200 
Set-Cookie: deeClientSeesionId=14443A044F8BBF4BDBA2131B92061CDA; Path=/; HttpOnly;deeI18n=zh_CN; Max-Age=86400000; Expires=Wed, 06-Sep-2028 04:08:51 GMT; Path=/
Content-Type: text/html;charset=UTF-8
Content-Length: 168
Date: Thu, 11 Dec 2025 04:08:51 GMT
==============================================

fmt.Printf("hearder %s\n", headers)

## 运行时候会报错
[ERR] Could not read nuclei-ignore file: open /Users/qq/Library/Application Support/nuclei/.nuclei-ignore: no such file or directory
搜关键词
.nuclei-ignore   会看到相关代码
./kscan -f 'body="sys/ui/extend/theme/default/style/icon.css"' -fs 20 -c
mac下需要复制一份 .nuclei-ignore 文件到如下目录， 因为他默认有个忽略模版，暂时还没想好怎么解决
/Users/qq/Library/Application Support/nuclei/.nuclei-ignore

此外可以屏蔽掉
from target list as found unresponsive 这里的输出信息 比如扫描的时候会输出
[INF] Skipped 60.190.96.156:80 from target list as found unresponsive 30 times
[INF] Skipped 60.190.96.156:80 from target list as found unresponsive 42 times
## 多条fofa语法合并
.\kscan -f "body='/weaver/ && body=/ecology '" -c -fs 10

## 这里有个问题就是 nacos这种 他无法匹配到
## poc 目录下不能有中文，否则无法加载poc

## 编译命令
$env:GOOS="linux"
$env:GOARCH="amd64"

$env:GOOS="windows"
$env:GOARCH="amd64"

zsh 终端
export GOOS=windows
export GOARCH=amd64
## 发现了写 yaml的时候
tags: sqli,致远互联-FE
tags可以使用中文，而且可以忽略大小写的

## fofa 查询后默认扫描的行为
sflag.BoolVar(&o.Check, "check", false) 在这里的 false控制

## 发现oa无法扫描到
http://211.143.78.190:7000
http://139.9.203.249:9088

## 2025.12.3 17:03
[+] 加载 20598 条指纹

## Joomla内容管理系统
是没有资产的，后面看看是否删除了 vps在跑脚本 去掉fofa上搜不到的资产
screen -r qwe 查看fofa的记录
ctrl+a d 退出
## 2025.12.3 发现 finger.json 里面的 keyword字段， 
程序处理的时候是严格大小写的，这里感觉可以优化

## bug 当扫描 tomcat rce的时候，发现了 存在问题，没有这个 tags
这里的想法是优化 finger字段，还有poc的 tags字段
http://123.58.224.8:59851/              Apache Tomcat;Apache-Tomcat;apache-tomcat
panic: 扫描执行失败: cause="No templates available"

goroutine 66 [running]:
kscan/nlscan.Nucleiscan({0xc013637e08, 0x1, 0x1}, {0xc007f34e10, 0x3, 0x3})
        D:/_笔记/_二开/kscan-master/nlscan/nucleiscan.go:114 +0x41b

## 2025.12.3 添加了漏扫
根据指纹匹配去漏扫，但是现在有个问题，就是指纹需要不断的优化，才能更好的扫描，因为tags目前没有做优化的
或者说 只用优化finger里面的字段 只要把finger.json 里面的字段优化为 tags里面的就好了
## 2025.11.24 日志 静态编译 finger.json

## 需要加一个poc的扫描
generateURLScanner() 这个函数输出了 匹配到的指纹

## 利用fofa 语法检索数据，并且 --check 用我们的指纹跑+探测存活 不会扫描端口
.\kscan -f header="Content-Length: 19" --check

## 对于带路径的 fofa语法需要
.\kscan -f "body='/jsoa/login.jsp'" --check -fs 100

fofa 搜索跳转页面的语法 是跳转后的页面特征，我们写的是跳转前的特征
## finger.json 维护的时候 
选择跳转之前的特征，否则抓不到

## 需要去掉颜色部分，因为颜色输出了东西 导致了捕获异常了
还有需要筛选一下 端口-指纹页面 的字段了
+ -Cn 参数即可 很多的爆破成功的看不到清楚
刷新问题，刷新频率慢一点或者手动刷新

## 如果需要 静态编译 finger.json
指纹修改文件:   lib/finger/finger.json
源码修改 lib/finger/finger.go

## 2025.10.23 修改
去除了 扫描addr的判断， 只保留了url和指纹信息

## 加载指纹修改为当前目录下 finger.json
这样方便修改指纹


