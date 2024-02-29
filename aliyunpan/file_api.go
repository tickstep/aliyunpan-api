package aliyunpan

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

	// FileUploadCheckPreHashParam 文件PreHash检测参数
	FileUploadCheckPreHashParam struct {
		// DriveId 网盘ID
		DriveId string `json:"drive_id"`
		// ParentFileId 父目录id，上传到根目录时填写 root
		ParentFileId string `json:"parent_file_id"`
		// Name 文件名称，按照 utf8 编码最长 1024 字节，不能以 / 结尾
		Name string `json:"name"`
		// Size 文件大小，单位为 byte。秒传必须
		Size int64 `json:"size"`
		// PreHash 针对大文件sha1计算非常耗时的情况， 可以先在读取文件的前1k的sha1， 如果前1k的sha1没有匹配的， 那么说明文件无法做秒传， 如果1ksha1有匹配再计算文件sha1进行秒传，这样有效边避免无效的sha1计算。
		PreHash string `json:"pre_hash"`
	}
)
