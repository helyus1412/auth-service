'use client';
import { Modal, Form, Input, message } from 'antd';
import axios from 'axios';

export default function UserFormModal({ open, onCancel, onSuccess, editingUser }: any) {
  const [form] = Form.useForm();

  const handleSubmit = async () => {
    const values = await form.validateFields();
    try {
      if (editingUser) {
        await axios.put(`http://localhost:8000/users/${editingUser.id}`, values);
        message.success('User updated');
      } else {
        await axios.post('http://localhost:8000/register', values);
        message.success('User added');
      }
      onSuccess();
    } catch {
      message.error('Failed to save user');
    }
  };

  return (
    <Modal
      open={open}
      title={editingUser ? 'Edit User' : 'Add User'}
      okText="Save"
      onCancel={onCancel}
      onOk={handleSubmit}
    >
      <Form form={form} layout="vertical" initialValues={editingUser || {}}>
        <Form.Item name="email" label="Email" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
      </Form>
    </Modal>
  );
}
