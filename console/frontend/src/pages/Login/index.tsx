import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { useState } from 'react';

import { login } from '@/services/auth';
import { useAuthStore } from '@/stores/auth';
import { message } from '@/utils/antdApp';
import { history } from '@/utils/history';

export default function LoginPage() {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      const res = await login(values);
      useAuthStore.getState().setTokens(res.token, res.refreshToken);
      message.success('登录成功');
      const params = new URLSearchParams(window.location.search);
      history.replace(params.get('redirect') || '/');
    } catch (error) {
      // 错误已由请求拦截器统一提示
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        padding: 16,
        boxSizing: 'border-box',
        background: 'linear-gradient(135deg, #1677ff 0%, #0b2545 100%)',
      }}
    >
      <div
        style={{
          width: '100%',
          maxWidth: 400,
          padding: 24,
          boxSizing: 'border-box',
          background: '#fff',
          borderRadius: 8,
          boxShadow: '0 8px 32px rgba(0,0,0,0.15)',
        }}
      >
        <LoginForm
          logo={<img alt="AdminBase" src="/logo.svg" style={{ width: 44, height: 44 }} />}
          title="AdminBase"
          subTitle="通用后台管理基座"
          loading={loading}
          containerStyle={{ padding: 0 }}
          contentStyle={{ minWidth: 0 }}
          onFinish={async (values: any) => {
            await handleSubmit(values);
          }}
        >
          <ProFormText
            name="username"
            fieldProps={{ size: 'large', prefix: <UserOutlined /> }}
            placeholder="用户名：admin"
            rules={[{ required: true, message: '请输入用户名' }]}
          />
          <ProFormText.Password
            name="password"
            fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
            placeholder="密码：123456"
            rules={[{ required: true, message: '请输入密码' }]}
          />
        </LoginForm>
      </div>
    </div>
  );
}
