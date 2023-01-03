# r8-retrace-web

鉴于混淆使用R8工具后，**proguardgui**工具无法完整解析出错误日志以及R8工具没有gui顾开发本工具

### 使用

1. [下载](https://github.com/Android-Mainli/r8-retrace-web/releases)

2. 确保已经安装android的`cmdline-tools`工具，并配置`ANDOIRD_HOME`环境变量

   或者直接配置r8的`retrace`命令至path

3. 启动本服务，如果是mac系统会直接打开浏览器 ，其他系统请查看启动服务后显示地址 选择一个进行访问

   如下：

   ```log
   地址:
           http://127.0.0.1:8082
   ```



> 不想使用默认端口可以添加-p <端号> 修改
>
> [R8 retrace官方文案](https://developer.android.google.cn/studio/command-line/retrace?hl=zh_cn)
