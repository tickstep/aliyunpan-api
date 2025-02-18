package aliyunpan

import "strings"

type (
	// FileList 文件列表
	FileList []*FileEntity

	// FileEntity 文件/文件夹信息
	FileEntity struct {
		// 网盘ID
		DriveId string `json:"driveId"`
		// 域ID
		DomainId string `json:"domainId"`
		// FileId 文件ID
		FileId string `json:"fileId"`
		// FileName 文件名
		FileName string `json:"fileName"`
		// FileSize 文件大小
		FileSize int64 `json:"fileSize"`
		// 文件类别 folder / file
		FileType string `json:"fileType"`
		// 创建时间
		CreatedAt string `json:"createdAt"`
		// 最后修改时间
		UpdatedAt string `json:"updatedAt"`
		// 后缀名，例如：dmg
		FileExtension string `json:"fileExtension"`
		// 文件上传ID
		UploadId string `json:"uploadId"`
		// 父文件夹ID
		ParentFileId string `json:"parentFileId"`
		// 内容CRC64校验值，只有文件才会有
		Crc64Hash string `json:"crc64Hash"`
		// 内容Hash值，只有文件才会有
		ContentHash string `json:"contentHash"`
		// 内容Hash计算方法，只有文件才会有，默认为：sha1
		ContentHashName string `json:"contentHashName"`
		// FilePath 文件的完整路径
		Path string `json:"path"`
		// Category 文件分类，例如：image/video/doc/others
		Category string `json:"category"`
		// SyncFlag 同步盘标记，该文件夹是否是同步盘的文件
		SyncFlag bool `json:"syncFlag"`
		// SyncMeta 如果是同步盘的文件夹，则这里会记录该文件对应的同步机器和目录等信息
		SyncMeta string `json:"syncMeta"`
		// Thumbnail 缩略图URL地址，只有相册文件才有
		Thumbnail string `json:"thumbnail"`
		// AlbumId 所属相册ID，只有相册文件才有
		AlbumId string `json:"albumId"`
	}

	FileOrderBy        string
	FileOrderDirection string
)

// IsFolder 是否是文件夹
func (f *FileEntity) IsFolder() bool {
	return f.FileType == "folder"
}

// IsFile 是否是文件
func (f *FileEntity) IsFile() bool {
	return f.FileType == "file"
}

// IsDriveRootFolder 是否是网盘根目录
func (f *FileEntity) IsDriveRootFolder() bool {
	return f.FileId == DefaultRootParentFileId
}

// 文件展示信息
func (f *FileEntity) String() string {
	builder := &strings.Builder{}
	builder.WriteString("文件ID: " + f.FileId + "\n")
	builder.WriteString("文件名: " + f.FileName + "\n")
	if f.IsFolder() {
		builder.WriteString("文件类型: 目录\n")
	} else {
		builder.WriteString("文件类型: 文件\n")
	}
	builder.WriteString("文件路径: " + f.Path + "\n")
	return builder.String()
}

// IsAlbumFile 是否是相册文件
func (f *FileEntity) IsAlbumFile() bool {
	return f.AlbumId != ""
}

// IsAlbumLivePhotoFile 是否是相册实况图片文件
func (f *FileEntity) IsAlbumLivePhotoFile() bool {
	return f.IsAlbumFile() && strings.HasSuffix(strings.ToLower(f.FileName), ".livp")
}

// TotalSize 获取目录下文件的总大小
func (fl FileList) TotalSize() int64 {
	var size int64
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		size += fl[k].FileSize
	}
	return size
}

// Count 获取文件总数和目录总数
func (fl FileList) Count() (fileN, directoryN int64) {
	for k := range fl {
		if fl[k] == nil {
			continue
		}

		if fl[k].IsFolder() {
			directoryN++
		} else {
			fileN++
		}
	}
	return
}

// ItemCount 文件数量
func (fl FileList) ItemCount() int {
	return len(fl)
}

// Item 文件项
func (fl FileList) Item(index int) *FileEntity {
	return fl[index]
}
