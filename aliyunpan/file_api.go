package aliyunpan

import (
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"time"
)

const (
	// DefaultRootParentFileId 网盘根目录默认ID
	DefaultRootParentFileId string = "root"

	FileOrderByName      FileOrderBy = "name"
	FileOrderByCreatedAt FileOrderBy = "created_at"
	FileOrderByUpdatedAt FileOrderBy = "updated_at"
	FileOrderBySize      FileOrderBy = "size"

	// FileOrderDirectionDesc 降序
	FileOrderDirectionDesc FileOrderDirection = "DESC"
	// FileOrderDirectionAsc 升序
	FileOrderDirectionAsc FileOrderDirection = "ASC"

	// MaxRequestRetryCount 最大重试次数（应对请求频繁的错误限制）
	MaxRequestRetryCount = int64(10)

	// IllegalDownloadUrlPrefix 资源被屏蔽，提示资源非法链接
	IllegalDownloadUrlPrefix = "https://pds-system-file.oss-cn-beijing.aliyuncs.com/illegal"

	// DefaultChunkSize 默认分片大小，512KB
	DefaultChunkSize = int64(524288)

	// MaxPartNum 最大分片数量大小
	MaxPartNum = 10000

	// DefaultZeroSizeFileContentHash 0KB文件默认的SHA1哈希值
	DefaultZeroSizeFileContentHash = "DA39A3EE5E6B4B0D3255BFEF95601890AFD80709"

	// ShellPatternCharacters 文件名\文件路径通配符字符串
	ShellPatternCharacters = "*?[]"

	// PathSeparator 路径分隔符
	PathSeparator = "/"
)

type (
	// FileCopyParam 文件复制参数
	FileCopyParam struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// ToParentFileId 目标目录ID、根目录为 root
		ToParentFileId string `json:"to_parent_file_id"`
	}

	// FileAsyncTaskResult 文件异步操作返回值
	FileAsyncTaskResult struct {
		// DriveId 网盘id
		DriveId string `json:"drive_id"`
		// FileId 文件ID
		FileId string `json:"file_id"`
		// AsyncTaskId 异步任务id。 如果返回为空字符串，表示直接移动成功。 如果返回非空字符串，表示需要经过异步处理。
		AsyncTaskId string `json:"async_task_id"`
	}

	// FileBatchActionParam 文件批量操作参数
	FileBatchActionParam struct {
		// 网盘ID
		DriveId string `json:"drive_id"`
		// 文件ID
		FileId string `json:"file_id"`
	}

	// FileBatchActionResult 文件批量操作返回值
	FileBatchActionResult struct {
		// 文件ID
		FileId string
		// 是否成功
		Success bool
	}

	// HandleFileDirectoryFunc 处理文件或目录的元信息, 返回值控制是否退出递归
	HandleFileDirectoryFunc func(depth int, fdPath string, fd *FileEntity, apierr *apierror.ApiError) bool

	// FileListParam 文件列表参数
	FileListParam struct {
		OrderBy        FileOrderBy        `json:"order_by"`
		OrderDirection FileOrderDirection `json:"order_direction"`
		DriveId        string             `json:"drive_id"`
		ParentFileId   string             `json:"parent_file_id"`
		Limit          int                `json:"limit"`
		// Marker 下一页参数
		Marker string `json:"marker"`
	}

	// FileListResult 文件列表返回值
	FileListResult struct {
		FileList FileList `json:"file_list"`
		// NextMarker 不为空代表还有下一页
		NextMarker string `json:"next_marker"`
	}

	// FileGetPathResult 文件路径详情信息结果
	FileGetPathResult struct {
		// 每一个item对应一个目录，最顶层的目录是root放在最后
		// 例如路径：/myphoto/photo2022/photo01，则对应顺序为item[0]={"photo01"}, item[1]={"photo2022"}, item[2]={"myphoto"}, item[3]={"root"}(只有root目录下的子文件夹才会有)
		Items []struct {
			Trashed      bool      `json:"trashed"`
			DriveId      string    `json:"drive_id"`
			FileId       string    `json:"file_id"`
			CreatedAt    time.Time `json:"created_at"`
			DomainId     string    `json:"domain_id"`
			EncryptMode  string    `json:"encrypt_mode"`
			Hidden       bool      `json:"hidden"`
			Name         string    `json:"name"`
			ParentFileId string    `json:"parent_file_id"`
			Starred      bool      `json:"starred"`
			Status       string    `json:"status"`
			Type         string    `json:"type"`
			UpdatedAt    string    `json:"updated_at"`
			UserMeta     string    `json:"user_meta"`
		} `json:"items"`
	}

	// FileMoveParam 文件移动参数
	FileMoveParam struct {
		// 源网盘ID
		DriveId string `json:"drive_id"`
		// 源文件ID
		FileId string `json:"file_id"`
		// 目标网盘ID
		ToDriveId string `json:"to_drive_id"`
		// 目标文件夹ID
		ToParentFileId string `json:"to_parent_file_id"`
	}

	// FileMoveResult 文件移动返回值
	FileMoveResult struct {
		// 文件ID
		FileId string
		// 是否成功
		Success bool
	}

	// VideoGetPreviewPlayInfoParam 视频信息参数
	VideoGetPreviewPlayInfoParam struct {
		DriveId string `json:"drive_id"`
		// FileId 视频文件ID
		FileId string `json:"file_id"`
	}

	// VideoGetPreviewPlayInfoResult 视频信息返回值
	VideoGetPreviewPlayInfoResult struct {
		DomainId             string `json:"domain_id"`
		DriveId              string `json:"drive_id"`
		FileId               string `json:"file_id"`
		VideoPreviewPlayInfo struct {
			Category string `json:"category"`
			Meta     struct {
				Duration            float64 `json:"duration"`
				Width               int     `json:"width"`
				Height              int     `json:"height"`
				LiveTranscodingMeta struct {
					TsSegment    int `json:"ts_segment"`
					TsTotalCount int `json:"ts_total_count"`
					TsPreCount   int `json:"ts_pre_count"`
				} `json:"live_transcoding_meta"`
			} `json:"meta"`
			LiveTranscodingTaskList []struct {
				TemplateId     string `json:"template_id"`
				TemplateName   string `json:"template_name"`
				TemplateWidth  int    `json:"template_width"`
				TemplateHeight int    `json:"template_height"`
				Status         string `json:"status"`
				Stage          string `json:"stage"`
				URL            string `json:"url"`
			} `json:"live_transcoding_task_list"`
		} `json:"video_preview_play_info"`
	}

	// MkdirResult 创建文件夹返回值
	MkdirResult struct {
		ParentFileId string `json:"parent_file_id"`
		Type         string `json:"type"`
		FileId       string `json:"file_id"`
		DomainId     string `json:"domain_id"`
		DriveId      string `json:"drive_id"`
		FileName     string `json:"file_name"`
		EncryptMode  string `json:"encrypt_mode"`
	}
)

// NewFileEntityForRootDir 创建根目录"/"的默认文件信息
func NewFileEntityForRootDir() *FileEntity {
	return &FileEntity{
		FileId:       DefaultRootParentFileId,
		FileType:     "folder",
		FileName:     "/",
		ParentFileId: "",
		Path:         "/",
	}
}
