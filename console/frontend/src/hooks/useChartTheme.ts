import { useSyncExternalStore } from 'react';

import { getSettings, subscribeSettings } from '@/theme';

/**
 * 图表主题：跟随 SettingDrawer 的明暗与主色设置
 */
export function useChartTheme() {
  const settings = useSyncExternalStore(subscribeSettings, getSettings);
  const isDark = settings.navTheme === 'realDark';
  return {
    isDark,
    g2Theme: isDark ? 'classicDark' : 'classic',
    colorPrimary: settings.colorPrimary || '#1677ff',
  };
}
