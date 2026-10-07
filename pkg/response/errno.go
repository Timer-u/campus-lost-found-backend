package response

import "net/http"

// Errno 业务错误：一个 Code 对应一个固定 Msg，文案与 docs/openapi.yaml 中的示例逐字一致
type Errno struct {
	Code int
	Msg  string
	HTTP int
}

func (e *Errno) Error() string {
	return e.Msg
}

// WithMsg 返回替换了提示信息的副本，用于同一错误码需要携带具体原因的场景
func (e *Errno) WithMsg(msg string) *Errno {
	return &Errno{Code: e.Code, Msg: msg, HTTP: e.HTTP}
}

// 业务错误码定义（新增错误码必须先在 docs/openapi.yaml 中补充）
var (
	ErrInvalidParams         = &Errno{Code: 10001, Msg: "参数错误", HTTP: http.StatusBadRequest}
	ErrUnauthorized          = &Errno{Code: 10002, Msg: "未登录或登录已过期", HTTP: http.StatusUnauthorized}
	ErrForbidden             = &Errno{Code: 10003, Msg: "无权操作该资源", HTTP: http.StatusForbidden}
	ErrAdminUndeletable      = &Errno{Code: 10004, Msg: "管理员账号不允许注销", HTTP: http.StatusForbidden}
	ErrUsernameExists        = &Errno{Code: 10007, Msg: "用户名已存在", HTTP: http.StatusConflict}
	ErrInvalidCredentials    = &Errno{Code: 10008, Msg: "账号或密码错误", HTTP: http.StatusUnauthorized}
	ErrAccountDisabled       = &Errno{Code: 10009, Msg: "账号已被禁用", HTTP: http.StatusForbidden}
	ErrNeedLostAdmin         = &Errno{Code: 10010, Msg: "需要管理员权限", HTTP: http.StatusForbidden}
	ErrNeedSystemAdmin       = &Errno{Code: 10011, Msg: "需要系统管理员权限", HTTP: http.StatusForbidden}
	ErrItemNotFound          = &Errno{Code: 12001, Msg: "物品不存在", HTTP: http.StatusNotFound}
	ErrItemStatusNotAllowed  = &Errno{Code: 12002, Msg: "物品当前状态不允许此操作", HTTP: http.StatusConflict}
	ErrClaimNotAllowed       = &Errno{Code: 13001, Msg: "当前物品不允许提交认领申请", HTTP: http.StatusConflict}
	ErrClaimNotFound         = &Errno{Code: 13002, Msg: "认领申请不存在", HTTP: http.StatusNotFound}
	ErrClaimStatusNotAllowed = &Errno{Code: 13003, Msg: "认领申请当前状态不允许此操作", HTTP: http.StatusConflict}
	ErrUserNotFound          = &Errno{Code: 14001, Msg: "用户不存在", HTTP: http.StatusNotFound}
	ErrAnnouncementNotFound  = &Errno{Code: 15001, Msg: "公告不存在", HTTP: http.StatusNotFound}
	ErrImageTooLarge         = &Errno{Code: 16001, Msg: "图片超过大小限制", HTTP: http.StatusRequestEntityTooLarge}

	// ErrInternal 文档未定义，作为 panic 恢复等服务器内部错误的兜底码
	ErrInternal = &Errno{Code: 10000, Msg: "服务器内部错误", HTTP: http.StatusInternalServerError}
)
