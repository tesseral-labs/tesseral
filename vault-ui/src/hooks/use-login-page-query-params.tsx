import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";

interface LoginPageQueryParams {
  relayedSessionState: string;
  redirectURI: string;
  returnRelayedSessionTokenAsQueryParam: boolean;
}

export function useLoginPageQueryParams(): LoginPageQueryParams {
  const [searchParams, setSearchParams] = useSearchParams();

  const [relayedSessionState, setRelayedSessionState] = useState(() => {
    return localStorage.getItem("relayed-session-state");
  });

  const [redirectURI, setRedirectURI] = useState(() => {
    return localStorage.getItem("redirect-uri");
  });

  const [
    returnRelayedSessionTokenAsQueryParam,
    setReturnRelayedSessionTokenAsQueryParam,
  ] = useState(() => {
    return localStorage.getItem("return-relayed-session-token-as-query-param");
  });

  useEffect(() => {
    const newParams = new URLSearchParams();

    if (searchParams.has("relayed-session-state")) {
      setRelayedSessionState(searchParams.get("relayed-session-state"));
      localStorage.setItem(
        "relayed-session-state",
        searchParams.get("relayed-session-state")!,
      );
    }

    if (searchParams.has("redirect-uri")) {
      setRedirectURI(searchParams.get("redirect-uri"));
      localStorage.setItem("redirect-uri", searchParams.get("redirect-uri")!);
    }

    if (searchParams.has("return-relayed-session-token-as-query-param")) {
      setReturnRelayedSessionTokenAsQueryParam(
        searchParams.get("return-relayed-session-token-as-query-param"),
      );
      localStorage.setItem(
        "return-relayed-session-token-as-query-param",
        searchParams.get("return-relayed-session-token-as-query-param")!,
      );
    }

    newParams.delete("relayed-session-state");
    newParams.delete("redirect-uri");
    newParams.delete("return-relayed-session-token-as-query-param");
    setSearchParams(newParams);
  }, [searchParams, setSearchParams]);

  return useMemo(() => {
    return {
      relayedSessionState: relayedSessionState || "",
      redirectURI: redirectURI || "",
      returnRelayedSessionTokenAsQueryParam:
        returnRelayedSessionTokenAsQueryParam === "1",
    } as LoginPageQueryParams;
  }, [relayedSessionState, redirectURI, returnRelayedSessionTokenAsQueryParam]);
}
