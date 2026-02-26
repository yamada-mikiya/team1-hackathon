import { Header } from '@/components/Header';
import './globals.css';
import { Providers } from './providers';

export const metadata = {
  title: 'Circle Blog',
  description: '大学サークル向けブログ',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="min-h-screen bg-gray-50">
        <Providers>
          <Header />
          {children}
        </Providers>
      </body>
    </html>
  );
}
