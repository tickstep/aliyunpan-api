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
)
