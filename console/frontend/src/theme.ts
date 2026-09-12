import type { ProLayoutProps } from '@ant-design/pro-components';

import defaultSettings from '@/config/defaultSettings';

export type LayoutSettings = ProLayoutProps & {
  pwa?: boolean;
  logo?: string;
  waterMark?: boolean;
};

const STORAGE_KEY = 'adminbase_layout_settings';

export function loadSettings(): LayoutSettings {
  if (typeof localStorage === 'undefined') {
    return { ...defaultSettings };
  }
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw
      ? { ...defaultSettings, ...(JSON.parse(raw) as Partial<LayoutSettings>) }
      : { ...defaultSettings };
  } catch {
    return { ...defaultSettings };
  }
}

let settings: LayoutSettings = loadSettings();
const listeners = new Set<() => void>();

export function getSettings(): LayoutSettings {
  return settings;
}

export function setSettings(next: Partial<LayoutSettings>): void {
  settings = { ...settings, ...next };
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
  } catch {
    // ignore quota / privacy mode errors
  }
  listeners.forEach((listener) => listener());
}

export function subscribeSettings(listener: () => void): () => void {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}
