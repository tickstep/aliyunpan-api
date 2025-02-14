package aliyunpan_open

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan_open/openapi"
)

type (
	ShareAlbumList []*aliyunpan.AlbumEntity
)

// ShareAlbumListGetAll 获取共享相册列表
func (p *OpenPanClient) ShareAlbumListGetAll() (ShareAlbumList, *apierror.ApiError) {
	retryTime := 0

RetryBegin:
	opParam := &openapi.ShareAlbumListParam{}
	if result, err := p.apiClient.ShareAlbumList(opParam); err == nil {
		shareAlbumList := ShareAlbumList{}
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
