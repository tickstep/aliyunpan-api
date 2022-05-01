# aliyunpan-api
GO语言封装的 aliyunpan 阿里云盘接口API。可以基于该接口库实现对阿里云盘的二次开发。

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/tickstep/aliyunpan-api?tab=doc)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/modern-go/concurrent/master/LICENSE)

# 快速使用

导入包
```
import "github.com/tickstep/aliyunpan-api/aliyunpan"
```

使用浏览器获取到的RefreshToken创建PanClient实例
```
	// get access token
	refreshToken := "f34b54eba1...706f389"
	webToken, err := aliyunpan.GetAccessTokenFromRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("get acccess token error")
		return
	}
	
	// pan client
	panClient := aliyunpan.NewPanClient(*webToken, aliyunpan.AppLoginToken{})
```

调用PanClient相关方法可以实现对阿里云盘的相关操作
```
	// get user info
	ui, err := panClient.GetUserInfo()
	if err != nil {
		fmt.Println("get user info error")
		return
	}
	fmt.Println("当前登录用户：" + ui.Nickname)

	// do some file operation
	fi, _ := panClient.FileInfoByPath(ui.FileDriveId, "/我的文档")
	fmt.Println("\n我的文档 信息：")
	fmt.Println(fi)
```