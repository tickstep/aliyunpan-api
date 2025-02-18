package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

// ShareAlbumListGetAll 获取共享相册列表
func (p *OpenPanClient) ShareAlbumListGetAll() (aliyunpan.ShareAlbumList, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.ShareAlbumListParam{}
	if result, err := p.apiClient.ShareAlbumList(opParam); err == nil {
		shareAlbumList := aliyunpan.ShareAlbumList{}
		for _, item := range result.Items {
			shareAlbumList = append(shareAlbumList, &aliyunpan.AlbumEntity{
				AlbumId:     item.SharedAlbumId,
				Name:        item.Name,
				Description: item.Description,
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   item.UpdatedAt,
				FileCount:   0,
				ImageCount:  0,
				VideoCount:  0,
				IsSharing:   false,
				Owner:       "",
			})
		}
		return shareAlbumList, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// ShareAlbumListFileGetAll 获取指定相簿下的所有文件列表
func (p *OpenPanClient) ShareAlbumListFileGetAll(param *aliyunpan.ShareAlbumListFileParam) (aliyunpan.FileList, *apierror.ApiError) {
	internalParam := &aliyunpan.ShareAlbumListFileParam{
		AlbumId:             param.AlbumId,
		Marker:              param.Marker,
		Limit:               param.Limit,
		ImageThumbnailWidth: param.ImageThumbnailWidth,
	}

	fileList := aliyunpan.FileList{}
	result, err := p.ShareAlbumListFile(internalParam)
	if err != nil || result == nil {
		return nil, err
	}
	fileList = append(fileList, result.FileList...)

	// more page?
	for len(result.NextMarker) > 0 {
		internalParam.Marker = result.NextMarker
		result, err = p.ShareAlbumListFile(internalParam)
		if err == nil && result != nil {
			fileList = append(fileList, result.FileList...)
		} else {
			break
		}
	}
	return fileList, nil
}

// ShareAlbumListFile 获取共享相册文件列表
func (p *OpenPanClient) ShareAlbumListFile(param *aliyunpan.ShareAlbumListFileParam) (*aliyunpan.FileListResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.ShareAlbumListFileParam{
		AlbumId:             param.AlbumId,
		OrderBy:             "joined_at",
		OrderDirection:      "DESC",
		Marker:              param.Marker,
		Limit:               param.Limit,
		ImageThumbnailWidth: param.ImageThumbnailWidth,
	}
	if opParam.Limit <= 0 {
		opParam.Limit = 50
	}
	if opParam.ImageThumbnailWidth <= 0 {
		opParam.ImageThumbnailWidth = 480
	}
	if result, err := p.apiClient.ShareAlbumListFile(opParam); err == nil {
		shareAlbumFileList := aliyunpan.FileList{}
		for _, item := range result.Items {
			shareAlbumFileList = append(shareAlbumFileList, &aliyunpan.FileEntity{
				DriveId:         item.DriveId,
				DomainId:        "",
				FileId:          item.FileId,
				FileName:        item.Name,
				FileSize:        item.Size,
				FileType:        item.Type,
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
				FileExtension:   item.FileExtension,
				UploadId:        "",
				ParentFileId:    item.ParentFileId,
				Crc64Hash:       "",
				ContentHash:     item.ContentHash,
				ContentHashName: item.ContentHashName,
				Path:            "",
				Category:        item.Category,
				SyncFlag:        false,
				SyncMeta:        "",
				Thumbnail:       item.Thumbnail,
				AlbumId:         item.AlbumId,
			})
		}
		return &aliyunpan.FileListResult{
			FileList:   shareAlbumFileList,
			NextMarker: result.NextMarker,
		}, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}

// ShareAlbumGetFileDownloadUrl 获取共享相册下文件下载地址
func (p *OpenPanClient) ShareAlbumGetFileDownloadUrl(param *aliyunpan.ShareAlbumGetFileUrlParam) (*aliyunpan.ShareAlbumGetFileUrlResult, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.ShareAlbumGetFileUrlParam{
		AlbumId: param.AlbumId,
		DriveId: param.DriveId,
		FileId:  param.FileId,
	}
	if result, err := p.apiClient.ShareAlbumGetFileDownloadUrl(opParam); err == nil {
		return result, nil
	} else {
		// handle common error
		if apiErrorHandleResp := p.HandleAliApiError(err, &retryTime); apiErrorHandleResp.NeedRetry {
			goto RetryBegin
		} else {
			return nil, apiErrorHandleResp.ApiErr
		}
	}
}
