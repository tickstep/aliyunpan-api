// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/library-go/jsonhelper"
	"os"
)

type (
	userpw struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
)

func objToJsonStr(v interface{}) string {
	r,_ := jsoniter.MarshalToString(v)
	return string(r)
}

func main() {
	configFile, err := os.OpenFile("userpw.txt", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("read user info error")
		return
	}
	defer configFile.Close()

	userpw := &userpw{}
	err = jsonhelper.UnmarshalData(configFile, userpw)
	if err != nil {
		fmt.Println("read user info error")
		return
	}

	// do login
	webToken, _ := aliyunpan.GetAccessTokenFromRefreshToken("3a4ec58d38de42d78f049691bbeab180")
	fmt.Println(objToJsonStr(webToken))

	// pan client
	panClient := aliyunpan.NewPanClient(*webToken, aliyunpan.AppLoginToken{})

	// user info
	fmt.Println(" ")
	ui,_ := panClient.GetUserInfo()
	fmt.Println(objToJsonStr(ui))

	// file list
	//fmt.Println(" ")
	//fl,_ := panClient.FileList(&aliyunpan.FileListParam{
	//	DriveId: ui.DefaultDriveId,
	//	ParentFileId: aliyunpan.DefaultRootParentFileId,
	//	Limit: 10,
	//})
	//fmt.Println(objToJsonStr(fl))

	// file info
	//fmt.Println(" ")
	//fi,_ := panClient.FileInfoById(ui.DefaultDriveId, "root")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoById(ui.DefaultDriveId, "60f406509fc5f3874bff4aa194dfc52011122859")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoById(ui.DefaultDriveId, "60f3c5b938e72352187e4c6da13879adf489267e")
	//fmt.Println(objToJsonStr(fi))

	// file info by path
	//fmt.Println("file info by path")
	//fi,_ = panClient.FileInfoByPath(ui.DefaultDriveId, "/")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoByPath(ui.DefaultDriveId, "/aliyunpan-api")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoByPath(ui.DefaultDriveId, "/aliyunpan-api/test.txt")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoByPath(ui.DefaultDriveId, "/aliyunpan-api/folder1/folder2/")
	//fmt.Println(objToJsonStr(fi))
	//fi,_ = panClient.FileInfoByPath(ui.DefaultDriveId, "/aliyunpan-api/folder1/folder2/folder3/test.zip")
	//fmt.Println(objToJsonStr(fi))

	// file list by path
	//fmt.Println("all file info in path")
	//fl1 := panClient.FilesDirectoriesRecurseList(ui.DefaultDriveId, "/", nil)
	//fmt.Println(objToJsonStr(fl1))
	//fl1 = panClient.FilesDirectoriesRecurseList(ui.DefaultDriveId, "/aliyunpan-api/folder1", nil)
	//fmt.Println(objToJsonStr(fl1))

	// mkdir
	fmt.Println("mkdir")
	mk,_ := panClient.Mkdir(ui.DefaultDriveId, "60f40b0feb801719ca9947478ee1f236bbdef9d8", "我的文件ABC123")
	fmt.Println(objToJsonStr(mk))
	mk,_ = panClient.MkdirByFullPath(ui.DefaultDriveId, "/aliyunpan-api/f1/f2/f3我的文件ABC123")
	fmt.Println(objToJsonStr(mk))
}
