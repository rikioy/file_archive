# file_archive

## 功能说明
1、把源文件夹下的文件(jpg, mp4...) ，读取拍摄日期后，按照 year/month-day/这样的结构存入到目标文件夹下。
2、程序会记录文件md5的值在"year/album.db"中，进行去重。

## 配置文件说明 config.ini
```
# 文件原始路径，不支持子目录。
src_path=../testpicture/src
# 目标路径
dst_path=./dst
# ffprobe程序路径，请自行安ffmpeg
ffprobe_exe=ffprobe
# 是否复制文件，true=src中文件复制到dst目录，默认为移动目录
copy_mode=false                

# 扩展类型名对应的处理程序可自行添
[type]
.jpg=exif
.mp4=mp4
```

## 使用说明
1. 根据实际情况配置好config.

### 批量处理文件
```
file_archive process
```

### 单独处理文件
> 在批量处理之后，如果没有检测到文件的创建日期，可以通过这个手动命令把文件加入到某个日期的文件夹下
```
 file_archive add -f 文件名 -d 日期字符串
```
1. 文件名，可以在failed.log中查看处理失败的文件名
2. 日期字符串，要把文件加入到哪个日期下，例如格式“20060102”会把文件加入到2006/01-02目录下
