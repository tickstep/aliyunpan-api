package aliyunpan

const (
	User        UserRole = "user"
	UnknownRole UserRole = "unknown"

	Enabled       UserStatus = "enable"
	UnknownStatus UserStatus = "unknown"
)

type (
	UserRole   string
	UserStatus string

	// UserInfo 用户信息
	UserInfo struct {
		// DomainId 域ID
		DomainId string `json:"domainId"`
		// FileDriveId 备份（文件）网盘ID
		FileDriveId string `json:"fileDriveId"`
		// SafeBoxDriveId 保险箱网盘ID
		SafeBoxDriveId string `json:"safeBoxDriveId"`
		// AlbumDriveId 相册网盘ID
		AlbumDriveId string `json:"albumDriveId"`
		// ResourceDriveId 资源库网盘ID
		ResourceDriveId string `json:"resourceDriveId"`
		// 用户UID
		UserId string `json:"userId"`
		// UserName 用户名
		UserName string `json:"userName"`
		// CreatedAt 创建时间
		CreatedAt string `json:"createdAt"`
		// Email 邮箱
		Email string `json:"email"`
		// Phone 手机
		Phone string `json:"phone"`
		// Role 角色，默认是user
		Role UserRole `json:"role"`
		// Status 是否被禁用，enable / disable
		Status UserStatus `json:"status"`
		// Nickname 昵称，如果没有设置则为空
		Nickname string `json:"nickname"`
		// TotalSize 网盘空间总大小
		TotalSize uint64 `json:"totalSize"`
		// UsedSize 网盘已使用空间大小
		UsedSize uint64 `json:"usedSize"`
		// ThirdPartyVip “三方权益包”是否生效
		ThirdPartyVip bool `json:"thirdPartyVip"`
		// ThirdPartyVipExpire “三方权益包”过期时间
		ThirdPartyVipExpire string `json:"thirdPartyVipExpire"`
	}
)
