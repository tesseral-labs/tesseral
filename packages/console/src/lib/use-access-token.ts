import { API_URL } from '@/config';
import { RefreshResponse } from '@/gen/tesseral/frontend/v1/frontend_pb';
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

export function useAccessToken(): string | undefined {
  const [accessToken, setAccessToken] = useLocalStorage('access_token');
  const [hasFailure, setHasFailure] = useState(false);
  const [refreshPending, setRefreshPending] = useState(false);

  if (!hasFailure && (!accessToken || shouldRefresh(accessToken))) {
    if (!refreshPending) {
      setRefreshPending(true);
      fetch(`${API_URL}/api/frontend/v1/refresh`, {
        body: JSON.stringify({}),
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then((response) => {
          if (response.ok) {
            response.json().then((res: RefreshResponse) => {
              setAccessToken(res.accessToken);
              setRefreshPending(false);
            });
          } else {
            setRefreshPending(false);
            setHasFailure(true);
          }
        })
        .catch((error) => {
          console.error(error);
          setRefreshPending(false);
          setHasFailure(true);
        });
    }
  }

  return accessToken ?? undefined;
}

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
