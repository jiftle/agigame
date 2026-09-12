import {
  App,
  Modal as staticModal,
  message as staticMessage,
  notification as staticNotification,
} from 'antd';

type AppFeedback = ReturnType<typeof App.useApp>;

let feedback: AppFeedback | null = null;

export function setAntdFeedback(instance: AppFeedback): void {
  feedback = instance;
}

export const message = new Proxy({} as AppFeedback['message'], {
  get(_target, prop) {
    const source = (feedback?.message ?? staticMessage) as unknown as Record<
      PropertyKey,
      unknown
    >;
    return source[prop];
  },
});

export const modal = new Proxy({} as AppFeedback['modal'], {
  get(_target, prop) {
    const source = (feedback?.modal ?? staticModal) as unknown as Record<
      PropertyKey,
      unknown
    >;
    return source[prop];
  },
});

export const notification = new Proxy({} as AppFeedback['notification'], {
  get(_target, prop) {
    const source = (feedback?.notification ?? staticNotification) as unknown as Record<
      PropertyKey,
      unknown
    >;
    return source[prop];
  },
});
