'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/providers/auth-provider';

export function useRequireAuth() {
  const { token, isReady } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isReady) {
      return;
    }
    if (token === null) {
      router.replace('/');
    }
  }, [isReady, router, token]);
}

export function useRedirectIfAuthenticated() {
  const { token, isReady } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isReady) {
      return;
    }
    if (token) {
      router.replace('/rooms');
    }
  }, [isReady, router, token]);
}
