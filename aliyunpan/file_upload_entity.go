package aliyunpan

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/tickstep/library-go/requester/rio"
	"io"
	"math"
	"math/big"
	"net/http"
)

type (
	// UploadFunc 上传文件处理函数
	UploadFunc func(httpMethod, fullUrl string, headers map[string]string) (resp *http.Response, err error)

	// FileUploadPartInfoParam 上传文件分片参数。从1开始，最大为 10000
	FileUploadPartInfoParam struct {
		PartNumber int `json:"part_number"`
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

	// CreateFileUploadParam 创建上传文件参数
	CreateFileUploadParam struct {
		Name         string `json:"name"`
		DriveId      string `json:"drive_id"`
		ParentFileId string `json:"parent_file_id"`
		Size         int64  `json:"size"`
		// 上传文件分片参数，最大为 10000
		PartInfoList []FileUploadPartInfoParam `json:"part_info_list"`
		ContentHash  string                    `json:"content_hash"`
		// 默认为 sha1。可选：sha1，none
		ContentHashName string `json:"content_hash_name"`
		// 默认为 file
		Type string `json:"type"`
		// 默认为 auto_rename。可选：overwrite-覆盖网盘同名文件，auto_rename-自动重命名，refuse-无需检测
		CheckNameMode string `json:"check_name_mode"`

		ProofCode    string `json:"proof_code"`
		ProofVersion string `json:"proof_version"`

		// 分片大小
		// 不进行json序列化
		BlockSize int64 `json:"-"`
		// LocalCreatedAt 本地创建时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalCreatedAt string `json:"-"`
		// LocalModifiedAt 本地修改时间，只对文件有效，格式yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		LocalModifiedAt string `json:"-"`
	}

	FileUploadPartInfoResult struct {
		PartNumber        int    `json:"part_number"`
		UploadURL         string `json:"upload_url"`
		InternalUploadURL string `json:"internal_upload_url"`
		ContentType       string `json:"content_type"`
	}

	// CreateFileUploadResult 创建上传文件返回值
	CreateFileUploadResult struct {
		ParentFileId string                     `json:"parent_file_id"`
		PartInfoList []FileUploadPartInfoResult `json:"part_info_list"`
		UploadId     string                     `json:"upload_id"`
		// RapidUpload 是否秒传。true-已秒传，false-没有秒传，需要手动上传
		RapidUpload bool   `json:"rapid_upload"`
		Type        string `json:"type"`
		FileId      string `json:"file_id"`
		DomainId    string `json:"domain_id"`
		DriveId     string `json:"drive_id"`
		// FileName 保存在网盘的名称，因为网盘会自动重命名同名的文件
		FileName    string `json:"file_name"`
		EncryptMode string `json:"encrypt_mode"`
		Location    string `json:"location"`
	}

	// GetUploadUrlParam 获取上传数据链接参数
	GetUploadUrlParam struct {
		DriveId      string                    `json:"drive_id"`
		FileId       string                    `json:"file_id"`
		PartInfoList []FileUploadPartInfoParam `json:"part_info_list"`
		UploadId     string                    `json:"upload_id"`
	}

	// GetUploadUrlResult 获取上传数据链接返回值
	GetUploadUrlResult struct {
		DomainId     string                     `json:"domain_id"`
		DriveId      string                     `json:"drive_id"`
		FileId       string                     `json:"file_id"`
		PartInfoList []FileUploadPartInfoResult `json:"part_info_list"`
		UploadId     string                     `json:"upload_id"`
		CreateAt     string                     `json:"create_at"`
	}

	FileUploadRange struct {
		// 起始值，包含
		Offset int64
		// 总上传长度
		Len int64
	}

	// FileUploadChunkData 文件上传数据块
	FileUploadChunkData struct {
		Reader    io.Reader
		ChunkSize int64

		hasReadCount int64
	}

	// CompleteUploadFileParam 提交上传文件传输完成参数
	CompleteUploadFileParam struct {
		DriveId  string `json:"drive_id"`
		FileId   string `json:"file_id"`
		UploadId string `json:"upload_id"`
	}

	CompleteUploadFileResult struct {
		DriveId         string `json:"drive_id"`
		DomainId        string `json:"domain_id"`
		FileId          string `json:"file_id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		Size            int64  `json:"size"`
		UploadId        string `json:"upload_id"`
		ParentFileId    string `json:"parent_file_id"`
		Crc64Hash       string `json:"crc64_hash"`
		ContentHash     string `json:"content_hash"`
		ContentHashName string `json:"content_hash_name"`
		CreatedAt       string `json:"created_at"`
	}
)

func (d *FileUploadChunkData) Read(p []byte) (n int, err error) {
	realReadCount := int64(0)
	var buf []byte = p
	needCopy := false
	if (d.hasReadCount + int64(len(p))) > d.ChunkSize {
		realReadCount = d.ChunkSize - d.hasReadCount
		buf = make([]byte, realReadCount)
		needCopy = true
	}

	n, err = d.Reader.Read(buf)
	if needCopy {
		copy(p, buf)
	}
	d.hasReadCount += int64(n)
	return n, err
}

func (d *FileUploadChunkData) Len() int64 {
	return d.ChunkSize
}

// CalcProofCode 计算文件上传防伪码
func CalcProofCode(accessToken string, reader rio.ReaderAtLen64, fileSize int64) string {
	if fileSize == 0 { // empty file
		return ""
	}

	md5w := md5.New()
	md5w.Write([]byte(accessToken))
	md5bytes := md5w.Sum(nil)
	hashCode := hex.EncodeToString(md5bytes)[0:16]
	hashInteger, _ := new(big.Int).SetString(hashCode, 16)

	z := big.NewInt(0)
	startPosInteger := big.NewInt(0)
	z.Div(hashInteger, big.NewInt(fileSize))
	startPosInteger.Sub(hashInteger, big.NewInt(z.Int64()*fileSize))
	startPos := startPosInteger.Int64()

	endPos := startPos + 8
	if endPos > fileSize {
		endPos = fileSize
	}

	// read byte from file
	readCount := endPos - startPos
	proofBytes := make([]byte, readCount)
	reader.ReadAt(proofBytes, startPos)

	// calc the base64 string for read bytes
	return base64.StdEncoding.EncodeToString(proofBytes)
}

// GenerateFileUploadPartInfoList 根据文件大小自动生成分片
func GenerateFileUploadPartInfoList(size int64) []FileUploadPartInfoParam {
	return GenerateFileUploadPartInfoListWithChunkSize(size, DefaultChunkSize)
}

// GenerateFileUploadPartInfoListWithChunkSize 根据文件大小和指定的分片大小自动生成分片
func GenerateFileUploadPartInfoListWithChunkSize(size, chunkSize int64) []FileUploadPartInfoParam {
	r := []FileUploadPartInfoParam{}
	if size <= chunkSize {
		r = append(r, FileUploadPartInfoParam{
			PartNumber: 1,
		})
	} else {
		pageSize := int(math.Ceil(float64(size) / float64(chunkSize)))
		for i := 1; i <= pageSize; i++ {
			r = append(r, FileUploadPartInfoParam{
				PartNumber: i,
			})
		}
	}
	return r
}
