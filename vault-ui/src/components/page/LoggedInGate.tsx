import { Code, ConnectError } from "@connectrpc/connect";
import { useMutation } from "@connectrpc/connect-query";
import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  useSyncExternalStore,
} from "react";
import React from "react";
import { Outlet, useNavigate } from "react-router";

import { refresh } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { parseAccessToken } from "@/lib/parse-access-token";

export function LoggedInGate() {
  const accessToken = useAccessTokenInternal();
  if (!accessToken) {
    return null;
  }
  return <Outlet />;
}

export function clearAccessToken() {
  localStorage.removeItem("access_token");
}

// the below is adapted from tesseral-sdk-react

function useAccessTokenInternal(): string | undefined {
  const [accessToken, setAccessToken] = useAccessTokenState();
  const parsedAccessToken = useMemo(() => {
    if (!accessToken) {
      return undefined;
    }
    return parseAccessToken(accessToken);
  }, [accessToken]);

  const now = useDebouncedNow(1000 * 2); // Re-check every 2 seconds
  const accessTokenIsLikelyValid = useMemo(() => {
    if (!parsedAccessToken || !parsedAccessToken.exp) {
      return false;
    }
    return parsedAccessToken.exp > now / 1000;
  }, [parsedAccessToken, now]);

  const { mutateAsync: refreshAsync } = useMutation(refresh);
  const navigate = useNavigate();

  useEffect(() => {
    if (accessTokenIsLikelyValid) {
      return;
    }

    async function refreshAccessToken() {
      try {
        const { accessToken } = await refreshAsync({});
        setAccessToken(accessToken!);
      } catch (e) {
        if (e instanceof ConnectError && e.code === Code.Unauthenticated) {
          navigate("/login");
          return;
        }

        throw e;
      }
    }

    void refreshAccessToken();
  }, [accessTokenIsLikelyValid, setAccessToken, refreshAsync, navigate]);

  if (accessTokenIsLikelyValid) {
    return accessToken!;
  }

  return;
}

function useDebouncedNow(updatePeriodMillis: number): number {
  const [now, setNow] = useState(Date.now());
  useEffect(() => {
    const interval = setInterval(() => setNow(Date.now()), updatePeriodMillis);
    return () => clearInterval(interval);
  }, [updatePeriodMillis]);
  return now;
}

function useAccessTokenState(): [
  string | null,
  (accessToken: string | null) => void,
] {
  return useLocalStorage(`access_token`);
}

function useLocalStorage(
  key: string,
): [string | null, (value: string | null) => void] {
  const store = useSyncExternalStore(
    // subscribe
    (callback) => {
      window.addEventListener("storage", callback);
      return () => window.removeEventListener("storage", callback);
    },
    // getSnapshot
    () => localStorage.getItem(key),
  );

  const setState = useCallback(
    (value: string | null) => {
      if (value === null) {
        localStorage.removeItem(key);
      } else {
        localStorage.setItem(key, value);
      }

      // setItem only dispatches on other tabs; we need to dispatch for our own
      // tab too
      window.dispatchEvent(new Event("storage"));
    },
    [key],
  );

  return [store, setState];
}
