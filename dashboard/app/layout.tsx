'use client';
import { ConfigProvider, App as AntdApp } from 'antd';
import 'antd/dist/reset.css';
import './globals.css';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <ConfigProvider
          theme={{
            token: {
              colorPrimary: '#1677ff',
              borderRadius: 6,
            },
          }}
        >
          <AntdApp>{children}</AntdApp>
        </ConfigProvider>
      </body>
    </html>
  );
}
