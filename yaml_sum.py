from pathlib import Path

def count_yaml_files(directory_path):
    # 创建 Path 对象
    path = Path(directory_path)
    
    # rglob 表示递归匹配 (Recursive glob)
    # 会匹配 .yaml 和 .yml 后缀
    yaml_files = list(path.rglob("*.yaml")) + list(path.rglob("*.yml"))
    
    return len(yaml_files)

# 使用示例
dir_to_scan = "./nlscan/poc/dahua"
count = count_yaml_files(dir_to_scan)
print(f"目录 '{dir_to_scan}' 及其子目录下共有 {count} 个 YAML 文件。")