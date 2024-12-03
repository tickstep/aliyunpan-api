package aliyunpan

type (
	// ShareCreateParam 创建分享
	ShareCreateParam struct {
		DriveId string `json:"drive_id"`
		// 分享密码，4个字符，为空代码公开分享
		SharePwd string `json:"share_pwd"`
		// 过期时间，为空代表永不过期。时间格式必须是这种：2021-07-23 09:22:19
		Expiration string   `json:"expiration"`
		FileIdList []string `json:"file_id_list"`
	}
	ShareEntity struct {
		Creator   string `json:"creator"`
		DriveId   string `json:"drive_id"`
		ShareId   string `json:"share_id"`
		ShareName string `json:"share_name"`
		// SharePwd 密码，为空代表没有密码
		SharePwd   string   `json:"share_pwd"`
		ShareUrl   string   `json:"share_url"`
		FileIdList []string `json:"file_id_list"`
		SaveCount  int      `json:"save_count"`
		// Expiration 过期时间，为空代表永不过期
		Expiration string `json:"expiration"`
		UpdatedAt  string `json:"updated_at"`
		CreatedAt  string `json:"created_at"`
		// forbidden-已违规，enabled-正常
		Status    string      `json:"status"`
		FirstFile *FileEntity `json:"first_file"`
	}

	// FastShareCreateParam 创建快传分享
	FastShareCreateParam struct {
		DriveId    string   `json:"drive_id"`
		FileIdList []string `json:"file_id_list"`
	}
	FastShareFileItem struct {
		DriveId string `json:"drive_id"`
		FileId  string `json:"file_id"`
	}
	// FastShareCreateResult 创建快传返回值
	FastShareCreateResult struct {
		Expiration    string              `json:"expiration"`
		Thumbnail     string              `json:"thumbnail"`
		ShareName     string              `json:"share_name"`
		ShareId       string              `json:"share_id"`
		ShareUrl      string              `json:"share_url"`
		DriveFileList []FastShareFileItem `json:"drive_file_list"`
		FullShareMsg  string              `json:"full_share_msg"`
		ShareTitle    string              `json:"share_title"`
		ShareSubtitle string              `json:"share_subtitle"`
		Expired       bool                `json:"expired"`
	}
)
