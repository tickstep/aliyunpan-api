package aliyunpan

import "net/http"

type (
	DownloadFuncCallback func(httpMethod, fullUrl string, headers map[string]string) (resp *http.Response, err error)

	// FileDownloadRange 分片。0-100,101-200,201-300...
	FileDownloadRange struct {
		// 起始值，包含
		Offset int64
		// 结束值，包含
		End int64
	}

	// GetFileDownloadUrlParam 获取文件下载链接
	GetFileDownloadUrlParam struct {
		DriveId   string `json:"drive_id"`
		FileId    string `json:"file_id"`
		ExpireSec int    `json:"expire_sec"`
	}

	// GetFileDownloadUrlResult 获取文件下载链接返回值
	GetFileDownloadUrlResult struct {
		Method      string `json:"method"`
		Url         string `json:"url"`
		InternalUrl string `json:"internal_url"`
		CdnUrl      string `json:"cdn_url"`
		Expiration  string `json:"expiration"`
		Size        int64  `json:"size"`
		Ratelimit   struct {
			PartSpeed int64 `json:"part_speed"`
			PartSize  int64 `json:"part_size"`
		} `json:"ratelimit"`
		Description string `json:"description"`
	}
)
