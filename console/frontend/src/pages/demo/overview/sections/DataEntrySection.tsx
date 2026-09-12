import { UploadOutlined } from '@ant-design/icons';
import {
  Button,
  Cascader,
  Checkbox,
  DatePicker,
  Flex,
  Input,
  InputNumber,
  Radio,
  Rate,
  Select,
  Slider,
  Space,
  Switch,
  TimePicker,
  Transfer,
  TreeSelect,
  Upload,
} from 'antd';
import type { TransferProps } from 'antd';
import { useState } from 'react';

import DemoBlock from '../DemoBlock';

const { TextArea, Password } = Input;

const selectOptions = [
  { label: '选项一', value: 1 },
  { label: '选项二', value: 2 },
  { label: '选项三', value: 3, disabled: true },
];

const treeData = [
  {
    title: '父节点',
    value: 'parent',
    children: [{ title: '子节点', value: 'child' }],
  },
];

const cascaderOptions = [
  {
    value: 'zhejiang',
    label: '浙江',
    children: [{ value: 'hangzhou', label: '杭州' }],
  },
];

const transferData = Array.from({ length: 8 }).map((_, i) => ({
  key: String(i),
  title: `条目 ${i + 1}`,
}));

export default function DataEntrySection() {
  const [targetKeys, setTargetKeys] = useState<TransferProps['targetKeys']>(['0', '1']);

  return (
    <>
      <DemoBlock title="Input 输入框">
        <Flex gap="middle" wrap>
          <Input placeholder="基础输入" style={{ width: 200 }} />
          <Input placeholder="带前缀" prefix="@" style={{ width: 200 }} />
          <Input allowClear placeholder="可清空" defaultValue="内容" style={{ width: 200 }} />
          <Input.Password placeholder="密码" style={{ width: 200 }} />
          <InputNumber placeholder="数字" min={0} max={100} style={{ width: 160 }} />
        </Flex>
        <TextArea rows={2} placeholder="多行文本" style={{ maxWidth: 640 }} />
        <Password placeholder="独立密码输入（解构自 Input）" style={{ width: 260 }} />
      </DemoBlock>

      <DemoBlock title="Select / TreeSelect / Cascader">
        <Flex gap="middle" wrap>
          <Select
            placeholder="单选"
            options={selectOptions}
            defaultValue={1}
            style={{ width: 200 }}
          />
          <Select
            mode="multiple"
            placeholder="多选"
            options={selectOptions}
            defaultValue={[1, 2]}
            style={{ width: 240 }}
          />
          <TreeSelect
            treeData={treeData}
            placeholder="树选择"
            defaultValue="child"
            style={{ width: 200 }}
          />
          <Cascader
            options={cascaderOptions}
            placeholder="级联选择"
            style={{ width: 220 }}
          />
        </Flex>
      </DemoBlock>

      <DemoBlock title="DatePicker / TimePicker">
        <Space wrap>
          <DatePicker placeholder="日期" />
          <DatePicker.RangePicker />
          <TimePicker placeholder="时间" />
        </Space>
      </DemoBlock>

      <DemoBlock title="Switch / Slider / Rate / Radio / Checkbox">
        <Space size="large" wrap>
          <Switch defaultChecked />
          <Switch checkedChildren="开" unCheckedChildren="关" defaultChecked />
          <Slider min={0} max={100} defaultValue={40} style={{ width: 200 }} />
          <Rate defaultValue={3} />
        </Space>
        <Radio.Group
          defaultValue="a"
          options={[
            { label: '选项 A', value: 'a' },
            { label: '选项 B', value: 'b' },
          ]}
        />
        <Checkbox.Group
          defaultValue={['apple']}
          options={[
            { label: '苹果', value: 'apple' },
            { label: '香蕉', value: 'banana' },
          ]}
        />
      </DemoBlock>

      <DemoBlock title="Transfer / Upload">
        <Transfer
          dataSource={transferData}
          targetKeys={targetKeys}
          onChange={(keys) => setTargetKeys(keys)}
          render={(item) => item.title}
        />
        <Upload beforeUpload={() => false} defaultFileList={[]}>
          <Button icon={<UploadOutlined />}>选择文件（不实际上传）</Button>
        </Upload>
      </DemoBlock>
    </>
  );
}
