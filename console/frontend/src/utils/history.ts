/**
 * 极薄的导航封装。
 * 由 <BrowserRouter> 内的 HistoryBridge 注入 react-router 的 navigate，
 * 使非组件模块（axios 拦截器、工具函数）也能发起路由跳转。
 */
type NavigateOptions = { replace?: boolean };
type Navigator = (to: string, options?: NavigateOptions) => void;

let navigator: Navigator | null = null;

export function setNavigator(fn: Navigator): void {
  navigator = fn;
}

export const history = {
  push: (to: string) => navigator?.(to),
  replace: (to: string) => navigator?.(to, { replace: true }),
  get location() {
    return window.location;
  },
};
