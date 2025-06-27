import { Code, ConnectError } from "@connectrpc/connect";
import { useMutation } from "@connectrpc/connect-query";
import React, { createContext, useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router";

import { DOGFOOD_PROJECT_ID } from "@/config";
import { refresh } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

import { parseAccessToken } from "./parse-access-token";

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
  const navigate = useNavigate();
  const [error, setError] = useState<unknown>();
  const [accessToken, setAccessToken] = useState(() => {
    return getCookie(`tesseral_${DOGFOOD_PROJECT_ID}_access_token`);
  });
  const { mutateAsync: refreshAsync } = useMutation(refresh);
  const { accessTokenLikelyValid } = useAccessTokenLikelyValid(
    accessToken ?? "",
  );

  // whenever the access token is invalid or near-expired, refresh it
  useEffect(() => {
    if (accessTokenLikelyValid) {
      return;
    }

    (async () => {
      try {
        const { accessToken } = await refreshAsync({});
        setAccessToken(accessToken);
      } catch (e) {
        if (e instanceof ConnectError && e.code === Code.Unauthenticated) {
          navigate("/login");
        } else {
          setError(`Unexpected response from /api/frontend/refresh: ${e}`);
        }
      }
    })();
  }, [accessTokenLikelyValid, navigate, refreshAsync]);

  if (error) {
    throw error;
  }

  return accessToken;
}

const ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS = 10 * 1000;

function useAccessTokenLikelyValid(accessToken: string): {
  accessTokenLikelyValid: boolean;
  exp: number;
} {
  const now = useDebouncedNow(2 * 1000); // re-check expiration every 2 seconds
  return useMemo(() => {
    if (!accessToken) {
      return {
        accessTokenLikelyValid: false,
        exp: 0,
      };
    }
    const parsedAccessToken = parseAccessToken(accessToken);
    return {
      accessTokenLikelyValid:
        parsedAccessToken.exp! * 1000 > now + ACCESS_TOKEN_EXPIRY_BUFFER_MILLIS,
      exp: parsedAccessToken.exp,
    };
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
    .split("; ")
    .find((row) => row.startsWith(key + "="))
    ?.split("=")[1];
}
