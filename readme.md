# File Explorer

## Feature
* 文件浏览
  * 从配置文件根目录(默认为/)开始，浏览子项目(文件和目录)
  * 文件列表每个项目最右侧插入ExtraInfo属性，包含评分、标签、评论
  * 当前目录选中项目后，点击对应按钮进入打开ExtraInfo编辑Tab
* 编辑ExtraInfo
  * 选中预览列表内项目后，编辑输入栏修改，预览列表同步显示修改后内容
  * 输入栏同步显示被选中项目聚合值，*表示选中项存在多个值
  * 提交时预览列表内所有项目的改动都会生效
  * ExtraInfo保存在同目录.fe_meta_data/extra_info.{timestamp}.yml
  * 预览列表选中项目可点击按钮还原历史保存内容
* 目录树
  * 从配置根目录开始，单击展开子目录
  * 双击目录，新建文件Tab打开目录
* 收藏夹
  * 文件Tab内 收藏/取消收藏 当前目录
  * 每个收藏目录也是一棵目录树，单击展开子目录，双击新Tab打开目录

## Run
* windows下执行build.cmd，mac下执行build.sh打包应用，运行
  * gui用fyne编写，理论上也可以修改-os参数打包成安卓/ios应用
* 配置文件
  * win: .exe所在目录/conf.yml
  * linux: ~/.fe_config/conf.yml

## Arch
```
file_explorer
├─common
│  ├─logger                 日志api
│  ├─model                  结构体
│  └─utils                  工具api
├─resource                  静态资源
│  ├─font                   字体
│  └─picture                图片
├─service                   服务层
├─test                      测试
└─view                      视图
    ├─components            文件浏览器组件
    │  ├─main_panel         窗口主面板
    │  │  ├─extra_info_tab  ExtraInfo Tab
    │  │  └─file_tab        文件Tab
    │  └─side_panel         窗口侧边栏
    ├─packed_widgets        fyne原生widget修改版
    └─store                 上下文数据
```

## TODO
* 组件尺寸优化
* 文件操作
  * 复制、移动
  * 批量重命名
  * 操作日志
* 文件运行
  * 打开默认文件浏览器
  * call第三方程序打开文件
    * 音乐播放器
  * 文本编辑
* 影音库
