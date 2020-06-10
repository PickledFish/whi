# whi
采用异步批量访问Url，取响应状态码，网页标题，访问失败错误信息
#使用方法

- --input 输入文件，文件中url一行一条
- --output 结果输入文件，输出格式为html，默认为result.html

```shell script
# cmd 
go run whi.go --input=test_target.txt --output=result.html

# build 
go build whi.go
whi --input=test_target.txt --output=result.html
```

# 结果展示
![result](https://raw.githubusercontent.com/PickledFish/whi/master/result.png)