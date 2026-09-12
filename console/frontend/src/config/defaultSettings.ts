import type { ProLayoutProps } from '@ant-design/pro-components';

const defaultSettings: ProLayoutProps & {
  logo?: string;
  waterMark?: boolean;
} = {
  navTheme: 'light',
  colorPrimary: '#1677ff',
  layout: 'mix',
  contentWidth: 'Fluid',
  fixedHeader: false,
  fixSiderbar: true,
  colorWeak: false,
  title: 'GoBoy 控制台',
  logo: '/logo.svg',
  iconfontUrl: '',
  waterMark: false,
  token: {},
};

export default defaultSettings;
