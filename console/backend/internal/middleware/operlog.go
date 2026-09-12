package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/model"
	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/ctxuser"
)

// OperLog 操作日志中间件（记录写操作）
func OperLog(r *ghttp.Request) {
	start := gtime.TimestampMilli()
	r.Middleware.Next()

	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return
	}
	path := r.URL.Path
	if strings.HasPrefix(path, "/api/v1/auth/") {
		return
	}

	operName := ""
	if user := ctxuser.Get(r.Context()); user != nil {
		operName = user.Username
	}

	param := sanitizeParam(r.GetBodyString())
	if len(param) > 2000 {
		param = param[:2000]
	}

	status := consts.StatusEnabled
	errMsg := ""
	if err := r.GetError(); err != nil {
		status = consts.StatusDisabled
		errMsg = err.Error()
	}

	_ = service.Log().OperLogCreate(r.Context(), &model.OperLogCreateInput{
		Title:         operTitle(r.Method, path),
		BusinessType:  operBusinessType(r.Method),
		Method:        path,
		RequestMethod: r.Method,
		OperName:      operName,
		OperUrl:       path,
		OperIp:        r.GetClientIp(),
		OperParam:     param,
		Status:        status,
		ErrorMsg:      errMsg,
		Cost:          int(gtime.TimestampMilli() - start),
	})
}

func operBusinessType(method string) string {
	switch method {
	case http.MethodPost:
		return "insert"
	case http.MethodPut:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "other"
	}
}

// sensitiveKeys 需要脱敏的请求字段（小写）
var sensitiveKeys = map[string]bool{
	"password":        true,
	"oldpassword":     true,
	"newpassword":     true,
	"confirmpassword": true,
	"token":           true,
	"accesstoken":     true,
	"refreshtoken":    true,
}

// sanitizeParam 对请求体中的敏感字段脱敏后返回
func sanitizeParam(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || (raw[0] != '{' && raw[0] != '[') {
		return raw
	}
	var data interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return raw
	}
	maskSensitive(data)
	if b, err := json.Marshal(data); err == nil {
		return string(b)
	}
	return raw
}

func maskSensitive(v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			if sensitiveKeys[strings.ToLower(k)] {
				val[k] = "******"
				continue
			}
			maskSensitive(child)
		}
	case []interface{}:
		for _, child := range val {
			maskSensitive(child)
		}
	}
}

func operTitle(method, path string) string {
	titles := map[string]string{
		"/api/v1/system/user":          "用户管理",
		"/api/v1/system/user/resetPwd": "重置用户密码",
		"/api/v1/system/user/status":   "用户状态",
		"/api/v1/system/role":          "角色管理",
		"/api/v1/system/menu":          "菜单管理",
		"/api/v1/system/dept":          "部门管理",
		"/api/v1/system/dict/type":     "字典类型",
		"/api/v1/system/dict/data":     "字典数据",
		"/api/v1/system/config":        "参数配置",
		"/api/v1/system/log/login":     "登录日志",
		"/api/v1/system/log/oper":      "操作日志",
	}
	if title, ok := titles[path]; ok {
		return title
	}
	return method + " " + path
}
