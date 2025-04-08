import { API_URL, DOGFOOD_PROJECT_ID } from '@/config';
import React, { createContext, useEffect, useMemo, useState } from 'react';

const Context = createContext<string | undefined>(undefined);

export function AccessTokenProvider({ children }: { children?: React.ReactNode }) {
  const accessToken = useAccessTokenInternal();
  return <Context.Provider value={accessToken}>{children}</Context.Provider>;
}

export function useAccessToken() {
  return React.useContext(Context);
}

const ACCESS_TOKEN_NAME = `tesseral_${DOGFOOD_PROJECT_ID}_access_token`;

function useAccessTokenInternal() {
  const accessToken = useMemo(() => {
    return document.cookie.split(';').find((row) => row.trim().startsWith(`${ACCESS_TOKEN_NAME}=`))?.split('=')[1];
  }, [document.cookie]);

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

      if (response.ok) {
        return;
      }

      if (response.status === 401) {
        return // refresh failed, caller's responsibility to redirect to /login
      }

      console.error('internal useAccessToken(): error fetching refresh', response.status, await response.text());
    })();
  }, [accessTokenLikelyValid]);

  return accessToken
}

const ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS = 10 * 1000;

function useAccessTokenLikelyValid(accessToken: string): boolean {
  const now = useDebouncedNow(10 * 1000); // re-check expiration every 10 seconds
  return useMemo(() => {
    if (!accessToken) {
      return false;
    }
    const parsedAccessToken = parseAccessToken(accessToken);
    return parsedAccessToken.exp * 1000 > now + ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS;
  }, [accessToken, now]);
}

function parseAccessToken(accessToken: string): any {
  const claimsPart = accessToken.split('.')[1];
  const decodedClaims = new TextDecoder().decode(Uint8Array.from(atob(claimsPart), (c) => c.charCodeAt(0)));
  return JSON.parse(decodedClaims);
}

function useDebouncedNow(updatePeriodMillis: number): number {
  const [now, setNow] = useState(Date.now());
  useEffect(() => {
    const interval = setInterval(() => setNow(Date.now()), updatePeriodMillis);
    return () => clearInterval(interval);
  }, [updatePeriodMillis]);
  return now;
}
