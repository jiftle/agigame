// Package logic 统一注册所有业务逻辑实现，通过空导入触发各组件的 init 注册。
package logic

import (
	_ "agigame/console/backend/internal/logic/auth"
	_ "agigame/console/backend/internal/logic/config"
	_ "agigame/console/backend/internal/logic/dept"
	_ "agigame/console/backend/internal/logic/dict"
	_ "agigame/console/backend/internal/logic/emu"
	_ "agigame/console/backend/internal/logic/log"
	_ "agigame/console/backend/internal/logic/menu"
	_ "agigame/console/backend/internal/logic/role"
	_ "agigame/console/backend/internal/logic/user"
)
