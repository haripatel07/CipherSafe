import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import { Toaster } from 'react-hot-toast';
import { AuthHydrator } from '@/components/AuthHydrator';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'CipherSafe',
  description: 'Your Self-Hosted Secrets Manager',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={inter.className}>
        {/* This component ensures our auth state is loaded on the client */}
        <AuthHydrator />
        
        {/* For notifications */}
        <Toaster position="bottom-right" />
        
        <div className="min-h-screen bg-background text-foreground">
          {children}
        </div>
      </body>
    </html>
  );
}