import os

# 设置要遍历的根文件夹路径
root_folder = r"Q:\go_project\mxshop\userop-web"  # 使用原始字符串

# 遍历文件夹和子文件夹
for root, dirs, files in os.walk(root_folder):
    for file in files:
        # 检查文件名是否包含目标字符串
        if "【加微信.赠送精品IT课程】" in file:
            # 创建新的文件名，将目标字符串去掉
            new_name = file.replace("【加微信.赠送精品IT课程】", "")
            
            # 拼接完整的文件路径
            old_path = os.path.join(root, file)
            new_path = os.path.join(root, new_name)
            
            # 重命名文件
            os.rename(old_path, new_path)
            print(f"文件 '{file}' 已重命名为 '{new_name}'")

print("操作完成！")
