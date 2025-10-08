'use client';
import { useState } from 'react';
import { Button, Form, Input, Card, message, Typography } from 'antd';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import axios from 'axios';

const { Title } = Typography;

export default function RegisterPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      await axios.post('http://localhost:8000/register', values);
      message.success('Registration successful! Please login.');
      router.push('/login');
    } catch (err: any) {
      console.error(err);
      message.error(err?.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <Card title={<Title level={4}>Register</Title>} className="auth-card">
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="email" label="Email" rules={[{ required: true }]}>
            <Input placeholder="Enter your email" />
          </Form.Item>
          <Form.Item name="password" label="Password" rules={[{ required: true, min: 6 }]}>
            <Input.Password placeholder="Enter your password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            Register
          </Button>
          <div style={{ marginTop: 16, textAlign: 'center' }}>
            Already have an account? <Link href="/login">Login</Link>
          </div>
        </Form>
      </Card>
    </div>
  );
}
