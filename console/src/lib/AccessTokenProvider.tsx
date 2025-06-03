import { API_URL, DOGFOOD_PROJECT_ID } from '@/config';
import React, { createContext, useEffect, useMemo, useState } from 'react';
import { parseAccessToken } from './parse-access-token';

const Context = createContext<string | undefined>(undefined);

export function AccessTokenProvider({
  children,
}: {
  children?: React.ReactNode;
}) {
  const accessToken = useAccessTokenInternal();
  return <Context.Provider value={accessToken}>{children}</Context.Provider>;
}

export function useAccessToken() {
  return React.useContext(Context);
}

function useAccessTokenInternal(): string | undefined {
  const [error, setError] = useState<unknown>();
  const [accessToken, setAccessToken] = useState(() => {
    return getCookie(`tesseral_${DOGFOOD_PROJECT_ID}_access_token`);
  });
  const accessTokenLikelyValid = useAccessTokenLikelyValid(accessToken ?? '');

  // whenever the access token is invalid or near-expired, refresh it
  useEffect(() => {
    if (accessTokenLikelyValid) {
      return;
    }

    (async () => {
      const response = await fetch(`${API_URL}/api/frontend/v1/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({}),
        credentials: 'include',
      });

      if (response.status === 401) {
        // our refresh token is no good
        window.location.href = `/login`;
        return;
      }

      if (!response.ok) {
        setError(
          `Unexpected response from /api/frontend/refresh: ${response.status}`,
        );
        return;
      }

      const { accessToken } = await response.json();
      if (!accessToken) {
        setError('No access token returned from /api/frontend/refresh');
        return;
      }
      setAccessToken(accessToken);
    })();
  }, [accessTokenLikelyValid]);

  if (error) {
    throw error;
  }

  return accessToken;
}

const ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS = 10 * 1000;

function useAccessTokenLikelyValid(accessToken: string): boolean {
  const now = useDebouncedNow(10 * 1000); // re-check expiration every 10 seconds
  return useMemo(() => {
    if (!accessToken) {
      return false;
    }
    const parsedAccessToken = parseAccessToken(accessToken);
    return (
      parsedAccessToken.exp! * 1000 > now + ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS
    );
  }, [accessToken, now]);
}

function useDebouncedNow(updatePeriodMillis: number): number {
  const [now, setNow] = useState(Date.now());
  useEffect(() => {
    const interval = setInterval(() => setNow(Date.now()), updatePeriodMillis);
    return () => clearInterval(interval);
  }, [updatePeriodMillis]);
  return now;
}

function getCookie(key: string): string | undefined {
  return document.cookie
    .split('; ')
    .find((row) => row.startsWith(key + '='))
    ?.split('=')[1];
}
