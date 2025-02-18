package aliyunpan

import "github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"

type (
	// AlbumEntity 相薄实体
	AlbumEntity struct {
		Owner       string `json:"owner"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AlbumId     string `json:"album_id"`
		FileCount   int    `json:"file_count"`
		ImageCount  int    `json:"image_count"`
		VideoCount  int    `json:"video_count"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		IsSharing   bool   `json:"is_sharing"`
	}

	ShareAlbumList []*AlbumEntity

	// ShareAlbumListFileParam 获取共享相册文件列表参数
	ShareAlbumListFileParam struct {
		// AlbumId 共享相册唯一ID
		AlbumId string `json:"albumId"`
		// Marker 分页标记
		Marker string `json:"marker"`
		// Limit 返回文件数量，默认50
		Limit int `json:"limit"`
		// ImageThumbnailWidth 生成的图片缩略图宽度，默认480px
		ImageThumbnailWidth int `json:"image_thumbnail_width"`
	}

	// ShareAlbumGetFileUrlParam 获取共享相册下文件下载地址参数
	ShareAlbumGetFileUrlParam struct {
		// AlbumId 共享相册唯一ID
		AlbumId string `json:"albumId"`
		// DriveId 文件所属drive
		DriveId string `json:"drive_id"`
		// FileId 文件id
		FileId string `json:"file_id"`
	}
	// ShareAlbumGetFileUrlResult 获取共享相册下文件下载地址返回值
	ShareAlbumGetFileUrlResult struct {
		// Url 下载地址，如果文件是livp时为空
		Url string `json:"url"`
		// StreamsUrl 文件是livp时有值
		StreamsUrl *ShareAlbumFileStreamUrlItem `json:"streams_url"`
		// Expiration 下载地址有效时间
		Expiration string `json:"expiration"`
		// Size 文件大小
		Size int64 `json:"size"`
		// ContentHash 文件哈希
		ContentHash string `json:"content_hash"`
	}
	// ShareAlbumFileStreamUrlItem livp文件下载流地址
	ShareAlbumFileStreamUrlItem struct {
		// Heic livp图片，与 jpeg 不会同时返回。
		Heic string `json:"heic"`
		// Jpeg jpeg图片，与 heic 不会同时返回。
		Jpeg string `json:"jpeg"`
		// Mov livp动画
		Mov string `json:"mov"`
	}
)

func (a *AlbumEntity) CreatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.CreatedAt)
}
func (a *AlbumEntity) UpdatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.UpdatedAt)
}
