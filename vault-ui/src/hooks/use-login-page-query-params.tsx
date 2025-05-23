import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";

interface LoginPageQueryParams {
  relayedSessionState: string;
  redirectURI: string;
  returnRelayedSessionTokenAsQueryParam: boolean;
}

export function useLoginPageQueryParams(): [LoginPageQueryParams, string] {
  const [searchParams, setSearchParams] = useSearchParams();

  const [relayedSessionState] = useState(() => {
    return searchParams.get("relayed-session-state");
  });

  const [redirectURI] = useState(() => {
    return searchParams.get("redirect-uri");
  });

  const [returnRelayedSessionTokenAsQueryParam] = useState(() => {
    return searchParams.get("return-relayed-session-token-as-query-param");
  });

  useEffect(() => {
    const newParams = new URLSearchParams();
    newParams.delete("relayed-session-state");
    newParams.delete("redirect-uri");
    newParams.delete("return-relayed-session-token-as-query-param");
    setSearchParams(newParams);
  }, [setSearchParams]);

  const state = useMemo(() => {
    return {
      relayedSessionState: relayedSessionState || "",
      redirectURI: redirectURI || "",
      returnRelayedSessionTokenAsQueryParam:
        returnRelayedSessionTokenAsQueryParam === "1",
    } as LoginPageQueryParams;
  }, [relayedSessionState, redirectURI, returnRelayedSessionTokenAsQueryParam]);

  const serialized = useMemo(() => {
    const search = new URLSearchParams();
    if (relayedSessionState) {
      search.set("relayed-session-state", relayedSessionState);
    }
    if (redirectURI) {
      search.set("redirect-uri", redirectURI);
    }
    if (returnRelayedSessionTokenAsQueryParam) {
      search.set(
        "return-relayed-session-token-as-query-param",
        returnRelayedSessionTokenAsQueryParam,
      );
    }

    if (search.size === 0) {
      return "";
    }
    return `?${search.toString()}`;
  }, [redirectURI, relayedSessionState, returnRelayedSessionTokenAsQueryParam]);

  return [state, serialized];
}
