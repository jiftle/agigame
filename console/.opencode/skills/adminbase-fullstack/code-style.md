# 前端代码规范（React + TypeScript + antd 6）

适用范围：`frontend/src/**`。项目已有约定（`AGENTS.md`、`README.md`）优先于本文件。

---

## 一、通用规范

### 1. 命名

- 组件文件：页面 `pages/<模块>/index.tsx`；组件 `PascalCase.tsx`（如 `ChartCard.tsx`）。
- 变量 / 函数：`camelCase`；hook：`useXxx`（`hooks/useXxx.ts`）。
- 常量：`UPPER_SNAKE_CASE`；类型 / 接口：`PascalCase`。
- 服务函数：`queryXxxList` / `getXxx` / `createXxx` / `updateXxx` / `deleteXxx`。
- 权限码：`system:<模块>:<list|add|edit|remove>`。

### 2. 组件

- 职责单一；表单、弹窗、列表拆分独立组件。
- props 必须显式类型（`interface Props`），不用 `any`。
- 组件默认导出；导出与文件名一致。

### 3. 接口调用

- API 统一封装到 `src/services/<模块>.ts`，**不在组件里直接写 axios**。
- 走 `src/api/client.ts` 的 `request`，返回值取后端 `{code,message,data}` 包体的 `data`。
- 响应类型用 `ApiResult<T>` / `PageResult<T>`（`services/types.ts`）；请求/响应都要有类型。
- 错误由拦截器统一 `message.error`，业务层通常无需重复提示。
- 异步请求覆盖 **loading / empty / error 三态**。

### 4. 状态管理

- 局部状态用组件内 `useState`。
- 服务端数据：优先 TanStack Query（`useCurrentUser` 范例）；**列表页沿用 `ProTable` 的 `request` 回调**，不要臆造所有页面都用 `useQuery`。
- 全局仅登录态用 Zustand（`stores/auth.ts`）；避免滥用全局状态。

### 5. 样式

- 优先 antd 组件与布局（`Row`/`Col`/`Flex`/`Space`/`ProCard`），颜色/间距用 `theme.useToken()`。
- **禁止硬编码颜色**：如 `#ff4d4f`、`rgba(0,0,0,.15)`。删除类操作链接用 `Button type="link" danger` 或 `token.colorError`。
- 全局样式放 `src/global.less`；**本项目未实际使用 Tailwind 工具类**，不要新增 Tailwind class（避免两套样式体系）。
- 不用 antd 静态 `message`/`Modal.confirm`，用 `@/utils/antdApp`。

### 6. TypeScript

- 避免 `any`，用具体类型或 `unknown` + 收窄。
- 接口请求 / 响应定义类型；函数参数与返回值标注类型。
- 状态 / 字典映射集中定义（`valueEnum` 或常量映射），禁止散落硬编码。

### 7. 代码质量

- 不留 `console.log` 调试代码、注释掉的代码、未使用的变量与 import。
- 复杂逻辑写必要注释（中文）；简单代码不写注释。
- 不引入新依赖（除非与用户确认）。

---

## 二、目录与文件结构

```
frontend/src/
├── api/
│   └── client.ts              # 唯一 axios 实例 + 拦截器（token / 401 刷新）
├── services/
│   ├── types.ts               # ApiResult<T> / PageResult<T> 等通用类型
│   └── <模块>.ts              # 接口封装 + 业务类型
├── pages/<模块>/
│   ├── index.tsx              # 页面：布局 + ProTable / ModalForm
│   └── components/            # 页面级子组件（按需）
├── components/                # 跨页面通用组件（ChartCard 等）
├── hooks/                     # useCurrentUser / useChartTheme 等
├── router/index.tsx           # 路由 + 权限守卫
├── config/menu.tsx            # 侧边栏菜单 + access
├── access.tsx                 # useAccess / <Access>
├── stores/                    # Zustand（登录态）
└── utils/                     # antdApp / history
```

拆分原则：

- 页面超过 ~300 行或将出现重复逻辑时，抽子组件 / hook。
- 列表页：表格列定义可留在页面内（项目现状），复杂时抽到同模块。
- 弹窗表单：简单作为子组件，复杂抽独立组件。

---

## 三、代码模板

### 3.1 services/<模块>.ts

```ts
import request from '@/api/client';

import type { ApiResult, PageResult } from './types';

export interface Xxx {
  id?: number;
  name?: string;
  status?: number;
}

export async function queryXxxList(params: {
  pageNum?: number;
  pageSize?: number;
  name?: string;
  status?: number;
}): Promise<PageResult<Xxx>> {
  const res = await request<ApiResult<PageResult<Xxx>>>('/system/xxx/list', { params });
  return res.data;
}

export async function getXxx(id: number): Promise<Xxx> {
  const res = await request<ApiResult<{ xxx: Xxx }>>(`/system/xxx/${id}`);
  return res.data.xxx;
}

export async function createXxx(data: Xxx): Promise<void> {
  await request('/system/xxx', { method: 'POST', data });
}

export async function updateXxx(data: Xxx): Promise<void> {
  await request('/system/xxx', { method: 'PUT', data });
}

export async function deleteXxx(ids: number[]): Promise<void> {
  await request('/system/xxx', { method: 'DELETE', data: { ids } });
}
```

### 3.2 列表页 index.tsx

```tsx
import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { ModalForm, PageContainer, ProFormText, ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm } from 'antd';
import { useRef, useState } from 'react';

import { Access, useAccess } from '@/access';
import type { Xxx } from '@/services/system';
import { createXxx, deleteXxx, queryXxxList, updateXxx } from '@/services/system';
import { message } from '@/utils/antdApp';

export default function XxxPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Xxx | undefined>();

  const reload = () => actionRef.current?.reload();

  const columns: ProColumns<Xxx>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '名称', dataIndex: 'name' },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: { 1: { text: '正常', status: 'Success' }, 0: { text: '停用', status: 'Error' } },
    },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:xxx:edit']}>
          <a
            onClick={() => {
              setEditing(record);
              setOpen(true);
            }}
          >
            编辑
          </a>
        </Access>,
        <Access key="remove" accessible={access['system:xxx:remove']}>
          <Popconfirm key="remove" title="确认删除？" onConfirm={() => deleteXxx([record.id as number]).then(reload)}>
            <Button type="link" danger size="small" style={{ padding: 0 }}>
              删除
            </Button>
          </Popconfirm>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<Xxx>
        headerTitle="Xxx 列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        form={{ syncToUrl: false }}
        request={async (params) => {
          const res = await queryXxxList({
            pageNum: params.current,
            pageSize: params.pageSize,
            name: params.name,
            status: params.status,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:xxx:add']}>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(undefined); setOpen(true); }}>
              新增
            </Button>
          </Access>,
        ]}
      />

      <ModalForm<Xxx>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑' : '新增'}
        open={open}
        onOpenChange={setOpen}
        width={520}
        initialValues={editing || { status: 1 }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateXxx({ ...values, id: editing.id });
          } else {
            await createXxx(values);
          }
          message.success(editing ? '修改成功' : '新增成功');
          reload();
          return true;
        }}
      >
        <ProFormText name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]} />
      </ModalForm>
    </PageContainer>
  );
}
```

### 3.3 详情 / 表单页

- 用 `PageContainer title="..."` 包裹；内容用 `ProCard` / `Card` + `Descriptions` / `ProForm`。
- 只读详情用 `Descriptions`；纵向流程用 `Steps` / `Timeline`（见 `pages/demo/detail`）。
- 表单用 `ProForm*` 组件，`onFinish` 返回 `true` 以关闭 loading。

---

## 四、关键规则（必须遵守）

- **权限**：页面用 `router` 的 `guard(...)`；按钮用 `useAccess()` + `<Access>`。
- **请求**：统一 `@/api/client`，服务层返回 `res.data`；不另建实例。
- **提示**：统一 `@/utils/antdApp` 的 `message`/`modal`，禁用 antd 静态方法。
- **颜色**：用 `theme.useToken()`，禁止硬编码色值。
- **状态/字典**：用 `valueEnum` 或集中映射，禁止硬编码。
- **三态**：loading / empty / error 必须覆盖。
- **不新增依赖**，不引入 Umi / Tailwind 工具类。
- **优先复用** `ProTable`、`ModalForm`、`ProCard`、`ChartCard`、`@/utils/antdApp` 等既有能力。

## 五、验证

```bash
cd frontend && pnpm tsc && pnpm build
```

类型正确性由 `tsc` 保证（项目无 ESLint / Prettier）。写 antd 代码先查后写：`antd info/demo/token`，改完 `antd lint <path>`。

改动 antd / 主题 / 全局样式后，按 SKILL.md「界面验证协议」过一遍 `/demo/overview` 与 `/demo/status`（明暗两态）。新增组件用法时同步补进组件总览页。
