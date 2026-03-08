import './globals.css';
import type { ReactNode } from 'react';
import { AuthProvider } from '@/app/providers/auth-provider';
import { ControlHeader } from '@/app/components/control-header';

export const metadata = {
  title: '✦ ITO ゲーム ✦',
  description: 'リアルタイムカードゲーム ITO'
};

export default function RootLayout({
  children
}: {
  children: ReactNode;
}) {
  return (
    <html lang="ja">
      <body>
        <AuthProvider>
          <main>
            <ControlHeader />
            {children}
          </main>
        </AuthProvider>
      </body>
    </html>
  );
}
