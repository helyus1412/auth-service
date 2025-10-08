'use client';
import { useState } from 'react';
import { Button, Form, Input, Card, message, Typography } from 'antd';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import axios from 'axios';

const { Title } = Typography;

export default function LoginPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const res = await axios.post('http://localhost:8000/login', values);
      localStorage.setItem('token', res.data.token);
      message.success('Login successful');
      router.push('/users');
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || 'Invalid email or password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <Card title={<Title level={4}>Login</Title>} className="auth-card">
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="email" label="Email" rules={[{ required: true }]}>
            <Input placeholder="Enter your email" />
          </Form.Item>
          <Form.Item name="password" label="Password" rules={[{ required: true }]}>
            <Input.Password placeholder="Enter your password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            Login
          </Button>
          <div style={{ marginTop: 16, textAlign: 'center' }}>
            Don't have an account? <Link href="/register">Register</Link>
          </div>
        </Form>
      </Card>
    </div>
  );
}
