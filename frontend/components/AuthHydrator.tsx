'use client';

import { useAuthStore } from '@/stores/authStore';
import { useEffect, useState } from 'react';

export function AuthHydrator() {
  const [isHydrated, setIsHydrated] = useState(false);
  const { isAuthenticated } = useAuthStore(); 

  useEffect(() => {
    setIsHydrated(true);
  }, []);

  return null; 
}