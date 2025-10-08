'use client';
import { useEffect, useState } from 'react';
import {
  App,
  Button,
  Table,
  Space,
  Typography,
} from 'antd';
import { useRouter } from 'next/navigation';
import axios from 'axios';
import UserFormModal from '../components/UserFormModal';

const { Title } = Typography;

interface User {
  id: number;
  email: string;
}

export default function UserPage() {
  const router = useRouter();
  const { message, modal } = App.useApp(); // ✅ context-aware modal + message
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await axios.get('http://localhost:8000/users');
      const data = Array.isArray(res.data) ? res.data : res.data.data || [];
      setUsers(data);
    } catch (err) {
      message.error('Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  const deleteUser = async (id: number) => {
    modal.confirm({
      title: 'Confirm Delete',
      content: 'Are you sure you want to delete this user?',
      okText: 'Yes, Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await axios.delete(`http://localhost:8000/users/${id}`);
          message.success('User deleted successfully');
          fetchUsers();
        } catch (err) {
          message.error('Failed to delete user');
        }
      },
    });
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const columns = [
    { title: 'ID', dataIndex: 'id' },
    { title: 'Email', dataIndex: 'email' },
    {
      title: 'Actions',
      render: (user: User) => (
        <Space>
          <Button
            type="link"
            onClick={() => {
              setEditingUser(user);
              setModalOpen(true);
            }}
          >
            Edit
          </Button>
          <Button
            type="link"
            danger
            onClick={() => deleteUser(user.id)}
          >
            Delete
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Title level={3}>User Management</Title>
      <Space style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          onClick={() => {
            setEditingUser(null);
            setModalOpen(true);
          }}
        >
          Add User
        </Button>
        <Button onClick={fetchUsers}>Refresh</Button>
        <Button
          onClick={() => {
            localStorage.clear();
            router.push('/login');
          }}
        >
          Logout
        </Button>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={users}
        loading={loading}
        bordered
      />

      {/* ✅ Add/Edit User Modal */}
      <UserFormModal
        open={isModalOpen}
        onCancel={() => setModalOpen(false)}
        onSuccess={() => {
          fetchUsers();
          setModalOpen(false);
        }}
        editingUser={editingUser}
      />
    </div>
  );
}
