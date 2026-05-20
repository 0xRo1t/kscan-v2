## 渗透测试自动化扫描器
kscan-v2

## 可维护性
可以对 finger.json进行维护， 以及 nlscan/poc 目录下的 nuclei-template模版进行编写，从而提高扫描器的漏洞探测能力

## why
本质上是对  https://github.com/0xRo1t/ezShell 以及 kscan的优化，同时更稳定

## 基础工作流
./kscan -f 'fofa语法' -fs fofa搜索条数 -c

./kscan -q 'quake语法' -qs quake搜索条数 -c

-c 是对 fofa/quake搜索的结果用我们自己的 finger.json里面的字典过一遍指纹，然后自动调用nuclei模版去扫描，已集成 nuclei官方sdk

## 日志
扫描后会自动生成日志文件，记录fofa/quake 搜索结果，以及扫描漏洞等详情日志，日志文件默认在当前目录下 scan_result/

## fofa/quake 
使用fofa或者quake 请在 kscan.go直接key := "" 替换掉对应值即可
