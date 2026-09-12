import {
  PageContainer,
  ProCard,
  ProForm,
  ProFormDatePicker,
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormSwitch,
  ProFormText,
  ProFormTextArea,
  StepsForm,
} from '@ant-design/pro-components';
import { Tabs } from 'antd';

import { message } from '@/utils/antdApp';

const ownerOptions = [
  { label: '张三', value: '张三' },
  { label: '李四', value: '李四' },
  { label: '王五', value: '王五' },
];

function BasicForm() {
  return (
    <ProForm
      style={{ maxWidth: 560 }}
      onFinish={async (values) => {
        // eslint-disable-next-line no-console
        console.log('基础表单:', values);
        message.success('提交成功');
        return true;
      }}
    >
      <ProFormText
        name="name"
        label="项目名称"
        rules={[{ required: true, message: '请输入项目名称' }]}
      />
      <ProFormSelect
        name="owner"
        label="负责人"
        options={ownerOptions}
        rules={[{ required: true, message: '请选择负责人' }]}
      />
      <ProFormDatePicker name="date" label="启动日期" />
      <ProFormDigit
        name="budget"
        label="预算(元)"
        min={0}
        fieldProps={{ precision: 0 }}
      />
      <ProFormRadio.Group
        name="level"
        label="优先级"
        initialValue="medium"
        options={[
          { label: '高', value: 'high' },
          { label: '中', value: 'medium' },
          { label: '低', value: 'low' },
        ]}
      />
      <ProFormSwitch name="notify" label="开启通知" initialValue />
      <ProFormTextArea name="desc" label="项目描述" />
    </ProForm>
  );
}

function StepForm() {
  return (
    <StepsForm
      onFinish={async (values) => {
        // eslint-disable-next-line no-console
        console.log('分步表单:', values);
        message.success('全部步骤已完成');
        return true;
      }}
    >
      <StepsForm.StepForm name="base" title="基本信息">
        <ProFormText
          name="name"
          label="项目名称"
          rules={[{ required: true, message: '请输入项目名称' }]}
        />
        <ProFormSelect name="owner" label="负责人" options={ownerOptions} />
      </StepsForm.StepForm>
      <StepsForm.StepForm name="config" title="参数配置">
        <ProFormDigit name="qps" label="QPS 上限" min={1} initialValue={100} />
        <ProFormSwitch name="cache" label="启用缓存" initialValue />
      </StepsForm.StepForm>
      <StepsForm.StepForm name="confirm" title="确认提交">
        <ProFormTextArea name="remark" label="备注" />
      </StepsForm.StepForm>
    </StepsForm>
  );
}

export default function DemoForm() {
  return (
    <PageContainer title="表单页">
      <ProCard>
        <Tabs
          items={[
            { key: 'base', label: '基础表单', children: <BasicForm /> },
            { key: 'step', label: '分步表单', children: <StepForm /> },
          ]}
        />
      </ProCard>
    </PageContainer>
  );
}
