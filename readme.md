# File Explorer

## Feature
### v0.1
* fix
  * ExtraInfo Tab Note输入栏显示多行
  * label换行显示修正
* 收藏夹
  * 打开收藏夹时检测收藏目录是否存在，存在失效目录时提示删除
    * 配置文件内提供 提示(0，默认)/直接删除(1)/忽略(2) 3种选项
* ExtraInfo
  * 新增url属性
    * 文件列表选中文件后，点击"打开url"按钮，默认浏览器打开被选文件urls
    * url支持文件名通配符"{name}"
* 文件浏览
  * 支持默认文件浏览器打开目录
  * 左上搜索框输入绝对路径打开目录
### v0.0
* 文件浏览
  * 配置文件可设置根目录，默认为/，双击进入子目录
* 目录树
  * 单击展开子目录
  * 双击打开目录
* 收藏夹
  * 文件Tab根据当前目录是否已收藏，工具栏显示"收藏/取消收藏"按钮
  * 每个收藏目录也是一棵目录树，单击展开子目录，双击打开目录
* ExtraInfo
  * 属性：评分、标签、评论
  * 文件ExtraInfo保存在同目录.fe_meta_data/extra_info.{timestamp}.yml内
  * 文件列表右侧显示ExtraInfo，双击标签、评论cell可显示全部内容
  * 文件列表选中项目后，点击"编辑ExtraInfo"按钮创建新Tab批量编辑
  * ExtraInfo编辑Tab
    * 预览列表选中项目后，编辑输入栏修改，预览列表同步显示修改后内容
    * 输入栏同步显示被选中项目聚合值，*表示选中项对应属性存在多个值
    * 点击"Cancel"按钮还原被选项目
    * 点击"Submit"按钮提交所有项目改动，提交时创建新yml文件存储ExtraInfo
    * 点击"还原"按钮显示目录下历史yml文件修改时间，点击历史文件后选中项目还原





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
  * 音乐播放器
  * 文本编辑
  * call第三方程序
* 影音库
