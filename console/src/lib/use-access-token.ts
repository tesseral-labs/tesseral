import { RefreshResponse } from '@/gen/tesseral/frontend/v1/frontend_pb';
import { useMutation } from '@tanstack/react-query';
import { useState } from 'react';

interface User {
  id: string;
  email: string;
}

export const useUser = (): User | undefined => {
  const accessToken = useAccessToken();
  if (!accessToken || accessToken.length === 0) {
    return;
  }

  const claims = JSON.parse(base64Decode(accessToken.split('.')[1]));
  return {
    id: claims.user.id,
    email: claims.user.email,
  };
};

export const useAccessToken = (): string | undefined => {
  const [hasFailure, setHasFailure] = useState(false);
  const [accessToken, setAccessToken] = useLocalStorage('access_token');
  const refresh = useMutation({
    mutationKey: ['refresh'],
    mutationFn: async () => {
      const response = await fetch(
        `https://auth.console.tesseral.example.com/api/frontend/v1/refresh`,
        {
          credentials: 'include',
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: '{}',
        },
      );

      if (!response.ok) {
        return;
      }

      return ((await response.json()) as RefreshResponse).accessToken;
    },
    retry: 0,
  });

  if (!hasFailure && (!accessToken || shouldRefresh(accessToken))) {
    if (!refresh.isPending) {
      refresh.mutate(undefined, {
        onError: () => {
          setHasFailure(true);
        },
        onSuccess: (accessToken: string) => {
          if (accessToken) {
            setHasFailure(false);
            setAccessToken(accessToken);
          }
        },
      });
    }
  }

  return accessToken ?? undefined;
};

// how far in advance of its expiration an access token gets refreshed
const ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS = 10;

const shouldRefresh = (accessToken: string): boolean => {
  const refreshAt =
    parseAccessTokenExpiration(accessToken) -
    ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS;
  const now = Math.floor(new Date().getTime() / 1000);
  return refreshAt < now;
};

const parseAccessTokenExpiration = (accessToken: string): number => {
  return JSON.parse(base64Decode(accessToken.split('.')[1])).exp;
};

const base64Decode = (s: string): string => {
  const binaryString = atob(s);

  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }

  return new TextDecoder().decode(bytes);
};

const useLocalStorage = (
  key: string,
): [string | null, (_: string | null) => void] => {
  const [value, setValue] = useState<string | null>(localStorage.getItem(key));

  return [
    value,
    (value) => {
      if (value === null) {
        localStorage.removeItem(key);
        setValue(null);
        return;
      }

      localStorage.setItem(key, value);
      setValue(value);
    },
  ];
};
