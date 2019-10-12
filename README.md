# file_archive

## 功能说明
1、把源文件夹下的文件(jpg, mp4...) ，读取拍摄日期后，按照 year/month-day/这样的结构存入到目标文件夹下。
2、程序会记录文件md5的值在"year/album.db"中，用于去重。

## 配置文件说明 config.ini
```
# 目标路径
dst_path=./dst

# 处理失败文件路径
fail_path=./failed

# ffprobe程序路径，请自行安ffmpeg
ffprobe_exe=ffprobe              

# 扩展类型名对应的处理程序可自行添
# 例如要添加对.bmp 格式支持,添加如下
# .bmp=exif
[type]
.jpg=exif
.JPG=exif
.mp4=mp4
```

## 使用说明
1. 根据实际情况配置好config.

### 批量处理文件
```
file_archive process -f 原始文件夹完整路径
```

### 单独处理文件
> 在批量处理之后，如果没有检测到文件的创建日期，可以通过这个手动命令把文件加入到某个日期的文件夹下
```
 file_archive add -f 文件路径 -d 日期字符串
```
1. 文件名，可以在failed.log中查看处理失败的文件名
2. 日期字符串，要把文件加入到哪个日期下，例如格式“20060102”会把文件加入到2006/01-02目录下

## 其他
### ffmpeg 下载
 * win https://ffmpeg.zeranoe.com/builds/
 * linux https://ffmpeg.org/download.html
### 识别格式
理论上有exif信息的图片类，或者ffprobe能读取的视频类文件都支持
