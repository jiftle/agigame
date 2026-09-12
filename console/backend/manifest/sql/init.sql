-- AdminBase 初始化脚本 (SQLite)
-- 说明: 所有语句幂等, 可重复执行

PRAGMA foreign_keys = ON;

-- ============================================================
-- 部门表
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_dept (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id   INTEGER NOT NULL DEFAULT 0,
    ancestors   TEXT    NOT NULL DEFAULT '',
    name        TEXT    NOT NULL DEFAULT '',
    leader      TEXT    NOT NULL DEFAULT '',
    phone       TEXT    NOT NULL DEFAULT '',
    email       TEXT    NOT NULL DEFAULT '',
    sort        INTEGER NOT NULL DEFAULT 0,
    status      INTEGER NOT NULL DEFAULT 1,
    remark      TEXT    NOT NULL DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);

-- ============================================================
-- 用户表
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_user (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    dept_id     INTEGER NOT NULL DEFAULT 0,
    username    TEXT    NOT NULL DEFAULT '',
    password    TEXT    NOT NULL DEFAULT '',
    nickname    TEXT    NOT NULL DEFAULT '',
    email       TEXT    NOT NULL DEFAULT '',
    phone       TEXT    NOT NULL DEFAULT '',
    sex         INTEGER NOT NULL DEFAULT 0,
    avatar      TEXT    NOT NULL DEFAULT '',
    status      INTEGER NOT NULL DEFAULT 1,
    login_ip    TEXT    NOT NULL DEFAULT '',
    login_date  DATETIME,
    remark      TEXT    NOT NULL DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_user_username ON sys_user (username) WHERE deleted_at IS NULL;

-- ============================================================
-- 角色表
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_role (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL DEFAULT '',
    code        TEXT    NOT NULL DEFAULT '',
    sort        INTEGER NOT NULL DEFAULT 0,
    data_scope  INTEGER NOT NULL DEFAULT 1,
    status      INTEGER NOT NULL DEFAULT 1,
    remark      TEXT    NOT NULL DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_role_code ON sys_role (code) WHERE deleted_at IS NULL;

-- ============================================================
-- 菜单/权限表
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_menu (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id   INTEGER NOT NULL DEFAULT 0,
    title       TEXT    NOT NULL DEFAULT '',
    name        TEXT    NOT NULL DEFAULT '',
    path        TEXT    NOT NULL DEFAULT '',
    component   TEXT    NOT NULL DEFAULT '',
    icon        TEXT    NOT NULL DEFAULT '',
    type        TEXT    NOT NULL DEFAULT 'C',
    perms       TEXT    NOT NULL DEFAULT '',
    sort        INTEGER NOT NULL DEFAULT 0,
    visible     INTEGER NOT NULL DEFAULT 1,
    status      INTEGER NOT NULL DEFAULT 1,
    redirect    TEXT    NOT NULL DEFAULT '',
    is_frame    INTEGER NOT NULL DEFAULT 0,
    is_cache    INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);

-- ============================================================
-- 用户-角色关联
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_user_role (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    role_id     INTEGER NOT NULL,
    created_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_user_role ON sys_user_role (user_id, role_id);

-- ============================================================
-- 角色-菜单关联
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_role_menu (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id     INTEGER NOT NULL,
    menu_id     INTEGER NOT NULL,
    created_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_role_menu ON sys_role_menu (role_id, menu_id);

-- ============================================================
-- 字典类型
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_dict_type (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL DEFAULT '',
    type        TEXT    NOT NULL DEFAULT '',
    status      INTEGER NOT NULL DEFAULT 1,
    remark      TEXT    NOT NULL DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_dict_type_type ON sys_dict_type (type) WHERE deleted_at IS NULL;

-- ============================================================
-- 字典数据
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_dict_data (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    dict_sort   INTEGER NOT NULL DEFAULT 0,
    dict_label  TEXT    NOT NULL DEFAULT '',
    dict_value  TEXT    NOT NULL DEFAULT '',
    dict_type   TEXT    NOT NULL DEFAULT '',
    is_default  INTEGER NOT NULL DEFAULT 0,
    status      INTEGER NOT NULL DEFAULT 1,
    remark      TEXT    NOT NULL DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);

-- ============================================================
-- 参数配置
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_config (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    config_name   TEXT    NOT NULL DEFAULT '',
    config_key    TEXT    NOT NULL DEFAULT '',
    config_value  TEXT    NOT NULL DEFAULT '',
    config_type   INTEGER NOT NULL DEFAULT 0,
    remark        TEXT    NOT NULL DEFAULT '',
    created_at    DATETIME,
    updated_at    DATETIME,
    deleted_at    DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_config_key ON sys_config (config_key) WHERE deleted_at IS NULL;

-- ============================================================
-- 登录日志
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_login_log (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username    TEXT    NOT NULL DEFAULT '',
    ip          TEXT    NOT NULL DEFAULT '',
    location    TEXT    NOT NULL DEFAULT '',
    browser     TEXT    NOT NULL DEFAULT '',
    os          TEXT    NOT NULL DEFAULT '',
    status      INTEGER NOT NULL DEFAULT 1,
    msg         TEXT    NOT NULL DEFAULT '',
    login_time  DATETIME
);

-- ============================================================
-- 操作日志
-- ============================================================
CREATE TABLE IF NOT EXISTS sys_oper_log (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    title           TEXT    NOT NULL DEFAULT '',
    business_type   TEXT    NOT NULL DEFAULT 'other',
    method          TEXT    NOT NULL DEFAULT '',
    request_method  TEXT    NOT NULL DEFAULT '',
    oper_name       TEXT    NOT NULL DEFAULT '',
    oper_url        TEXT    NOT NULL DEFAULT '',
    oper_ip         TEXT    NOT NULL DEFAULT '',
    oper_param      TEXT    NOT NULL DEFAULT '',
    json_result     TEXT    NOT NULL DEFAULT '',
    status          INTEGER NOT NULL DEFAULT 1,
    error_msg       TEXT    NOT NULL DEFAULT '',
    cost            INTEGER NOT NULL DEFAULT 0,
    oper_time       DATETIME
);

-- ============================================================
-- 种子数据
-- ============================================================

-- 部门
INSERT OR IGNORE INTO sys_dept (id, parent_id, ancestors, name, leader, phone, email, sort, status, remark, created_at, updated_at)
VALUES (1, 0, '0', '飞鱼科技', '管理员', '13800000000', 'admin@adminbase.dev', 1, 1, '总部', datetime('now','localtime'), datetime('now','localtime'));

-- 用户 (密码均为 123456)
INSERT OR IGNORE INTO sys_user (id, dept_id, username, password, nickname, email, phone, sex, status, remark, created_at, updated_at)
VALUES (1, 1, 'admin', '$2a$10$paP5N3UvXnd9EjtQMLB2yeN69Zm65hZKQlOf.Y1SUPFBJiCRnxx1O', '超级管理员', 'admin@adminbase.dev', '13800000000', 1, 1, '内置超级管理员', datetime('now','localtime'), datetime('now','localtime'));

-- 角色
INSERT OR IGNORE INTO sys_role (id, name, code, sort, data_scope, status, remark, created_at, updated_at)
VALUES (1, '超级管理员', 'admin', 1, 1, 1, '拥有全部权限', datetime('now','localtime'), datetime('now','localtime'));
INSERT OR IGNORE INTO sys_role (id, name, code, sort, data_scope, status, remark, created_at, updated_at)
VALUES (2, '普通用户', 'common', 2, 5, 1, '默认只读角色', datetime('now','localtime'), datetime('now','localtime'));

-- 用户-角色
INSERT OR IGNORE INTO sys_user_role (user_id, role_id, created_at) VALUES (1, 1, datetime('now','localtime'));

-- 菜单
INSERT OR IGNORE INTO sys_menu (id, parent_id, title, name, path, component, icon, type, perms, sort, visible, status, redirect, is_frame, is_cache, created_at, updated_at) VALUES
(100, 0,   '仪表盘',   'Dashboard',    '/dashboard',       './Dashboard',            'DashboardOutlined', 'C', '',                 1, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(1,   0,   '系统管理', 'System',       '/system',          '',                       'SettingOutlined',   'M', '',                 2, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(2,   1,   '用户管理', 'SystemUser',   '/system/user',     './system/user',          'UserOutlined',      'C', 'system:user:list',   1, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(3,   1,   '角色管理', 'SystemRole',   '/system/role',     './system/role',          'TeamOutlined',      'C', 'system:role:list',   2, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(4,   1,   '菜单管理', 'SystemMenu',   '/system/menu',     './system/menu',          'MenuOutlined',      'C', 'system:menu:list',   3, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(5,   1,   '部门管理', 'SystemDept',   '/system/dept',     './system/dept',          'ApartmentOutlined', 'C', 'system:dept:list',   4, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(6,   1,   '字典管理', 'SystemDict',   '/system/dict',     './system/dict',          'BookOutlined',      'C', 'system:dict:list',   5, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(7,   1,   '参数配置', 'SystemConfig', '/system/config',   './system/config',        'ControlOutlined',   'C', 'system:config:list', 6, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(8,   1,   '操作日志', 'SystemOperLog','/system/oper-log', './system/oper-log',      'FileTextOutlined',  'C', 'system:operlog:list',7, 1, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(9,   1,   '登录日志', 'SystemLoginLog','/system/login-log','./system/login-log',    'LoginOutlined',     'C', 'system:loginlog:list',8,1,1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime'));

-- 按钮权限
INSERT OR IGNORE INTO sys_menu (id, parent_id, title, name, path, component, icon, type, perms, sort, visible, status, redirect, is_frame, is_cache, created_at, updated_at) VALUES
(201, 2, '用户查询', 'SystemUserQuery',  '', '', '', 'F', 'system:user:query',   1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(202, 2, '用户新增', 'SystemUserAdd',    '', '', '', 'F', 'system:user:add',     2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(203, 2, '用户修改', 'SystemUserEdit',   '', '', '', 'F', 'system:user:edit',    3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(204, 2, '用户删除', 'SystemUserRemove', '', '', '', 'F', 'system:user:remove',  4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(205, 2, '重置密码', 'SystemUserReset',  '', '', '', 'F', 'system:user:resetPwd',5, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(301, 3, '角色查询', 'SystemRoleQuery',  '', '', '', 'F', 'system:role:query',   1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(302, 3, '角色新增', 'SystemRoleAdd',    '', '', '', 'F', 'system:role:add',     2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(303, 3, '角色修改', 'SystemRoleEdit',   '', '', '', 'F', 'system:role:edit',    3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(304, 3, '角色删除', 'SystemRoleRemove', '', '', '', 'F', 'system:role:remove',  4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(401, 4, '菜单查询', 'SystemMenuQuery',  '', '', '', 'F', 'system:menu:query',   1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(402, 4, '菜单新增', 'SystemMenuAdd',    '', '', '', 'F', 'system:menu:add',     2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(403, 4, '菜单修改', 'SystemMenuEdit',   '', '', '', 'F', 'system:menu:edit',    3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(404, 4, '菜单删除', 'SystemMenuRemove', '', '', '', 'F', 'system:menu:remove',  4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(501, 5, '部门查询', 'SystemDeptQuery',  '', '', '', 'F', 'system:dept:query',   1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(502, 5, '部门新增', 'SystemDeptAdd',    '', '', '', 'F', 'system:dept:add',     2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(503, 5, '部门修改', 'SystemDeptEdit',   '', '', '', 'F', 'system:dept:edit',    3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(504, 5, '部门删除', 'SystemDeptRemove', '', '', '', 'F', 'system:dept:remove',  4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(601, 6, '字典查询', 'SystemDictQuery',  '', '', '', 'F', 'system:dict:query',   1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(602, 6, '字典新增', 'SystemDictAdd',    '', '', '', 'F', 'system:dict:add',     2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(603, 6, '字典修改', 'SystemDictEdit',   '', '', '', 'F', 'system:dict:edit',    3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(604, 6, '字典删除', 'SystemDictRemove', '', '', '', 'F', 'system:dict:remove',  4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(701, 7, '参数查询', 'SystemConfigQuery', '', '', '', 'F', 'system:config:query', 1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(702, 7, '参数新增', 'SystemConfigAdd',   '', '', '', 'F', 'system:config:add',   2, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(703, 7, '参数修改', 'SystemConfigEdit',  '', '', '', 'F', 'system:config:edit',  3, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(704, 7, '参数删除', 'SystemConfigRemove','', '', '', 'F', 'system:config:remove',4, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(801, 8, '日志删除', 'SystemOperLogRemove', '', '', '', 'F', 'system:operlog:remove', 1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime')),
(901, 9, '日志删除', 'SystemLoginLogRemove','', '', '', 'F', 'system:loginlog:remove',1, 0, 1, '', 0, 1, datetime('now','localtime'), datetime('now','localtime'));

-- 超管拥有全部菜单
INSERT OR IGNORE INTO sys_role_menu (role_id, menu_id, created_at)
SELECT 1, id, datetime('now','localtime') FROM sys_menu;

-- 普通用户仅仪表盘
INSERT OR IGNORE INTO sys_role_menu (role_id, menu_id, created_at) VALUES (2, 100, datetime('now','localtime'));

-- 字典类型
INSERT OR IGNORE INTO sys_dict_type (id, name, type, status, remark, created_at, updated_at) VALUES
(1, '用户性别', 'sys_user_sex',    1, '用户性别列表', datetime('now','localtime'), datetime('now','localtime')),
(2, '通用状态', 'sys_normal_disable', 1, '启用/停用状态', datetime('now','localtime'), datetime('now','localtime'));

-- 字典数据
INSERT OR IGNORE INTO sys_dict_data (id, dict_sort, dict_label, dict_value, dict_type, is_default, status, remark, created_at, updated_at) VALUES
(1, 1, '男',   '1', 'sys_user_sex', 1, 1, '性别男', datetime('now','localtime'), datetime('now','localtime')),
(2, 2, '女',   '2', 'sys_user_sex', 0, 1, '性别女', datetime('now','localtime'), datetime('now','localtime')),
(3, 3, '未知', '0', 'sys_user_sex', 0, 1, '性别未知', datetime('now','localtime'), datetime('now','localtime')),
(4, 1, '正常', '1', 'sys_normal_disable', 1, 1, '正常状态', datetime('now','localtime'), datetime('now','localtime')),
(5, 2, '停用', '0', 'sys_normal_disable', 0, 1, '停用状态', datetime('now','localtime'), datetime('now','localtime'));

-- 参数配置
INSERT OR IGNORE INTO sys_config (id, config_name, config_key, config_value, config_type, remark, created_at, updated_at) VALUES
(1, '系统名称',   'sys.site.name',        'AdminBase 管理基座', 0, '系统显示名称', datetime('now','localtime'), datetime('now','localtime')),
(2, '默认密码',   'sys.user.initPwd',     '123456',             1, '新建用户默认密码', datetime('now','localtime'), datetime('now','localtime')),
(3, '登录验证码', 'sys.login.captcha',    'false',              0, '是否开启登录验证码', datetime('now','localtime'), datetime('now','localtime'));
