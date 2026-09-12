export interface MenuNode {
  id: number;
  parentId: number;
  title: string;
  name?: string;
  path?: string;
  component?: string;
  icon?: string;
  type: string;
  perms?: string;
  sort?: number;
  visible?: number;
  redirect?: string;
  isFrame?: number;
  isCache?: number;
  children?: MenuNode[];
}

export interface CurrentUser {
  id: number;
  username: string;
  nickname: string;
  avatar?: string;
  email?: string;
  phone?: string;
  sex?: number;
  deptId?: number;
  remark?: string;
}

export interface ApiResult<T> {
  code: number;
  message: string;
  data: T;
}

export interface PageResult<T> {
  total: number;
  list: T[];
}
