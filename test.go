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
		UserName     string `json:"username"`
		Password     string `json:"password"`
		RefreshToken string `json:"refreshToken"`
	}
)

func objToJsonStr(v interface{}) string {
	r, _ := jsoniter.MarshalToString(v)
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
	webToken, _ := aliyunpan.GetAccessTokenFromRefreshToken(userpw.RefreshToken)
	fmt.Println(objToJsonStr(webToken))

	// pan client
	panClient := aliyunpan.NewPanClient(*webToken, aliyunpan.AppLoginToken{})

	//// user info
	//fmt.Println(" ")
	//ui, _ := panClient.GetUserInfo()
	//fmt.Println(objToJsonStr(ui))
	//
	//// file list
	//fmt.Println(" ")
	//fl, _ := panClient.FileList(&aliyunpan.FileListParam{
	//	//OrderBy:        aliyunpan.FileOrderByName,
	//	//OrderDirection: aliyunpan.FileOrderDirectionAsc,
	//	DriveId:      ui.FileDriveId,
	//	ParentFileId: "610dfd8ab42d8eae886c4776927dca2a12dccb6a",
	//	Limit:        10,
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
	//fmt.Println("mkdir")
	//mk,_ := panClient.Mkdir(ui.DefaultDriveId, "60f40b0feb801719ca9947478ee1f236bbdef9d8", "我的文件ABC123")
	//fmt.Println(objToJsonStr(mk))
	//mk,_ = panClient.MkdirByFullPath(ui.DefaultDriveId, "/aliyunpan-api/f1/f2/f3我的文件ABC123")
	//fmt.Println(objToJsonStr(mk))

	// rename
	//fmt.Println("rename")
	//bb,_ := panClient.Rename(ui.DefaultDriveId, "60f432a37717e8d190cb443084d61409be1e44bc", "我的文件ABC123-1")
	//fmt.Println(bb)
	//bb,_ = panClient.Rename(ui.DefaultDriveId, "60f40b4794573a3eeb4b4e05904cebfc35328732", "ok-1.dmg")
	//fmt.Println(bb)

	// download url
	//fmt.Println("download url")
	//dp := &aliyunpan.GetFileDownloadUrlParam{
	//	DriveId: ui.FileDriveId,
	//	FileId:  "60f3c5b938e72352187e4c6da13879adf489267e",
	//}
	//gfdr, _ := panClient.GetFileDownloadUrl(dp)
	//fmt.Println(objToJsonStr(gfdr))
	//fmt.Println(gfdr.Url)

	// batch task
	//requests := aliyunpan.BatchRequestList{}
	//requests = append(requests, &aliyunpan.BatchRequest{
	//	Id:      "60bc44fcafaac4e737d14c969899d1ca553a7fa8",
	//	Method:  "POST",
	//	Url:     "/file/move",
	//	Headers: map[string]string{
	//		"Content-Type": "application/json",
	//	},
	//	Body:    map[string]interface{}{
	//		"drive_id": "19519221",
	//		"file_id": "60bc44fcafaac4e737d14c969899d1ca553a7fa8",
	//		"to_drive_id": "19519221",
	//		"to_parent_file_id": "60f61cf4f15322f69a4c4d4fb58bbcf8188e788b",
	//	},
	//})
	//batchParam := aliyunpan.BatchRequestParam{
	//	Requests: requests,
	//	Resource: "file",
	//}
	//result,_ := panClient.BatchTask("https://api.aliyundrive.com/v3/batch", &batchParam)
	//fmt.Println(objToJsonStr(result))

	// file move
	//fmp := []*aliyunpan.FileMoveParam{}
	//fmp = append(fmp, &aliyunpan.FileMoveParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44fcafaac4e737d14c969899d1ca553a7fa8",
	//	ToDriveId:      "19519221",
	//	ToParentFileId: "60f61cf4f15322f69a4c4d4fb58bbcf8188e788b",
	//})
	//fmp = append(fmp, &aliyunpan.FileMoveParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44fb9cc604dab65c4709ba27f2cbc8954a20",
	//	ToDriveId:      "19519221",
	//	ToParentFileId: "60f61cf4f15322f69a4c4d4fb58bbcf8188e788b",
	//})
	//result,_ := panClient.FileMove(fmp)
	//fmt.Println(objToJsonStr(result))

	// file delete
	//fdp := []*aliyunpan.FileBatchActionParam{}
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44f8b120740fb5534534845ccbb9b973c8c1",
	//})
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44fb7d68cc31cb024b19a7babe936ecb9af8",
	//})
	//result,_ := panClient.FileDelete(fdp)
	//fmt.Println(objToJsonStr(result))

	// recycle bin file delete
	//fdp := []*aliyunpan.FileBatchActionParam{}
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44f8367b636659564c2b9f057e8f5654241f",
	//})
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44f703d4b148194447639d04632674c6b465",
	//})
	//result,_ := panClient.RecycleBinFileDelete(fdp)
	//fmt.Println(objToJsonStr(result))

	// file starred
	//fdp := []*aliyunpan.FileBatchActionParam{}
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60bc44fb7d68cc31cb024b19a7babe936ecb9af8",
	//})
	//fdp = append(fdp, &aliyunpan.FileBatchActionParam{
	//	DriveId:        "19519221",
	//	FileId:         "60f61cf4f15322f69a4c4d4fb58bbcf8188e788b",
	//})
	//result,_ := panClient.FileStarred(fdp)
	//fmt.Println(objToJsonStr(result))

	// share list
	//slr,_ := panClient.ShareLinkList(ui.UserId)
	//fmt.Println(objToJsonStr(slr))

	// share cancel
	//links := []string{}
	//links = append(links, "EPPnhbvacC7")
	//links = append(links, "LSzZy6SFmGg")
	//slc,_ := panClient.ShareLinkCancel(links)
	//fmt.Println(objToJsonStr(slc))

	// share create
	//fileIdList := []string{}
	//fileIdList = append(fileIdList, "60f565ffe840c531d96d45cd9fca67a1a2528831")
	//fileIdList = append(fileIdList, "60bc44fb7d68cc31cb024b19a7babe936ecb9af8")
	//slc,_ := panClient.ShareLinkCreate(aliyunpan.ShareCreateParam{
	//	DriveId: "19519221",
	//	SharePwd: "a123",
	//	Expiration: "2021-07-30 07:18:07",
	//	FileIdList: fileIdList,
	//})
	//fmt.Println(objToJsonStr(slc))

	//albumList, _ := panClient.AlbumList(&aliyunpan.AlbumListParam{
	//	OrderBy:        aliyunpan.AlbumOrderByFileCount,
	//	OrderDirection: aliyunpan.AlbumOrderDirectionDesc,
	//	Limit:          2,
	//	Marker:         "WyI0ZDAwMWQ0ODU2NGY0M2IzYmM1NjYyODc0ZjA0YmJlNiIsIm4iLDAsImJqMjkjYzJhM2MyNmNmMDVjNDMxYmFkMzJjZDE3NmIyOTRmNjY2MjY1MTcyZiJd",
	//})
	//fmt.Println(objToJsonStr(albumList))

	r, _ := panClient.AlbumEdit(&aliyunpan.AlbumEditParam{
		AlbumId:     "7caa2962f1ea4da398c4696ab70e7a636202512b",
		Description: "古代诗歌",
		Name:        "663333",
	})
	fmt.Println(objToJsonStr(r))
}
