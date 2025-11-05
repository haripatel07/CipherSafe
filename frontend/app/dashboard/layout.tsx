'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/stores/authStore';
import { useRouter } from 'next/navigation';
import { Loader2, Lock } from 'lucide-react';
import toast from 'react-hot-toast';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, logout, token } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.replace('/login');
    }
  }, [isAuthenticated, router]);

  const handleLogout = () => {
    logout();
    toast.success('Logged out');
    router.push('/login');
  };

  // Show a loading screen while auth state is being determined
  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="animate-spin h-12 w-12" />
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col">
      <nav className="flex items-center justify-between p-4 border-b border-gray-800 bg-gray-900">
        <div className="flex items-center gap-2">
          <Lock className="text-blue-500" />
          <h1 className="text-xl font-bold">CipherSafe</h1>
        </div>
        <button
          onClick={handleLogout}
          className="px-4 py-2 text-sm font-medium bg-red-600 rounded-md hover:bg-red-700"
        >
          Logout
        </button>
      </nav>
      <main className="flex-grow p-4 md:p-8">
        {children}
      </main>
    </div>
  );
}