import { UserOutlined } from '@ant-design/icons';
import { PageContainer, ProCard, ProForm, ProFormText } from '@ant-design/pro-components';
import { useQueryClient } from '@tanstack/react-query';
import { Avatar, Descriptions, Form, Input, Tabs, Tag } from 'antd';
import { useSearchParams } from 'react-router-dom';

import { useCurrentUser } from '@/hooks/useCurrentUser';
import { changePassword, updateProfile, type UpdateProfileParams } from '@/services/auth';
import { useAuthStore } from '@/stores/auth';
import { message } from '@/utils/antdApp';
import { history } from '@/utils/history';

const AVATAR_RULE = /^https?:\/\/.+/i;
const PHONE_RULE = /^1[3-9]\d{9}$/;

function BasicSettings() {
  const { data } = useCurrentUser();
  const user = data?.user;
  const queryClient = useQueryClient();

  return (
    <ProCard>
      <div style={{ marginBottom: 24 }}>
        <Avatar size={72} src={user?.avatar || undefined} icon={<UserOutlined />} />
      </div>
      <ProForm<UpdateProfileParams>
        style={{ maxWidth: 480 }}
        initialValues={{
          nickname: user?.nickname,
          email: user?.email,
          phone: user?.phone,
          avatar: user?.avatar,
        }}
        onFinish={async (values) => {
          await updateProfile(values);
          message.success('保存成功');
          queryClient.invalidateQueries({ queryKey: ['auth', 'user'] });
          return true;
        }}
      >
        <ProFormText
          name="nickname"
          label="昵称"
          rules={[{ required: true, message: '请输入昵称' }]}
        />
        <ProFormText
          name="email"
          label="邮箱"
          rules={[{ type: 'email', message: '邮箱格式不正确' }]}
        />
        <ProFormText
          name="phone"
          label="手机号"
          rules={[{ pattern: PHONE_RULE, message: '手机号格式不正确' }]}
        />
        <ProFormText
          name="avatar"
          label="头像地址"
          placeholder="填写图片 URL，留空则使用默认图标"
          rules={[{ pattern: AVATAR_RULE, message: '请输入 http(s) 开头的图片地址' }]}
        />
      </ProForm>
    </ProCard>
  );
}

function SecuritySettings() {
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  return (
    <ProCard>
      <Form
        form={form}
        layout="vertical"
        style={{ maxWidth: 420 }}
        onFinish={async (values: {
          oldPassword: string;
          newPassword: string;
          confirmPassword: string;
        }) => {
          await changePassword({
            oldPassword: values.oldPassword,
            newPassword: values.newPassword,
          });
          message.success('密码修改成功，请重新登录');
          useAuthStore.getState().clear();
          queryClient.clear();
          history.replace('/login');
        }}
      >
        <Form.Item
          name="oldPassword"
          label="原密码"
          rules={[{ required: true, message: '请输入原密码' }]}
        >
          <Input.Password placeholder="请输入原密码" autoComplete="current-password" />
        </Form.Item>
        <Form.Item
          name="newPassword"
          label="新密码"
          rules={[
            { required: true, message: '请输入新密码' },
            { min: 6, max: 64, message: '密码长度需为 6-64 位' },
          ]}
        >
          <Input.Password placeholder="6-64 位" autoComplete="new-password" />
        </Form.Item>
        <Form.Item
          name="confirmPassword"
          label="确认新密码"
          dependencies={['newPassword']}
          rules={[
            { required: true, message: '请再次输入新密码' },
            ({ getFieldValue }) => ({
              validator(_rule, value) {
                if (!value || getFieldValue('newPassword') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('两次输入的密码不一致'));
              },
            }),
          ]}
        >
          <Input.Password placeholder="请再次输入新密码" autoComplete="new-password" />
        </Form.Item>
      </Form>
    </ProCard>
  );
}

function AccountInfo() {
  const { data } = useCurrentUser();
  const user = data?.user;
  const roles = data?.roles ?? [];
  const perms = data?.perms ?? [];
  const menus = data?.menus ?? [];
  const isSuper = perms.includes('*:*:*');

  return (
    <ProCard>
      <Descriptions column={1} bordered size="small">
        <Descriptions.Item label="用户名">{user?.username || '-'}</Descriptions.Item>
        <Descriptions.Item label="昵称">{user?.nickname || '-'}</Descriptions.Item>
        <Descriptions.Item label="角色">
          {roles.length ? roles.map((role) => <Tag key={role}>{role}</Tag>) : '-'}
        </Descriptions.Item>
        <Descriptions.Item label="权限数">{isSuper ? '全部' : perms.length}</Descriptions.Item>
        <Descriptions.Item label="菜单数">{menus.length}</Descriptions.Item>
      </Descriptions>
    </ProCard>
  );
}

export default function AccountPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const tab = searchParams.get('tab') ?? 'basic';

  return (
    <PageContainer title="个人中心">
      <ProCard>
        <Tabs
          activeKey={tab}
          onChange={(key) =>
            setSearchParams(key === 'basic' ? {} : { tab: key }, { replace: true })
          }
          items={[
            { key: 'basic', label: '基本资料', children: <BasicSettings /> },
            { key: 'security', label: '安全设置', children: <SecuritySettings /> },
            { key: 'about', label: '账号信息', children: <AccountInfo /> },
          ]}
        />
      </ProCard>
    </PageContainer>
  );
}
