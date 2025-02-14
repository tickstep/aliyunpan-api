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
)

func (a *AlbumEntity) CreatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.CreatedAt)
}
func (a *AlbumEntity) UpdatedAtStr() string {
	return apiutil.UnixTime2LocalFormat(a.UpdatedAt)
}
