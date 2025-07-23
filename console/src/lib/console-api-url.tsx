import React, {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

const Context = createContext<string | undefined>(undefined);

export function ConsoleApiUrlProvider({ children }: PropsWithChildren) {
  const [apiUrl, setApiUrl] = useState<string | undefined>(undefined);

  useEffect(() => {
    async function fetchConfig() {
      try {
        const response = await fetch("/config.json");
        if (response.ok) {
          const config = await response.json();
          setApiUrl(config.CONSOLE_API_URL);
        }
      } catch (error) {
        console.error("Failed to fetch config.json:", error);
      }
    }

    fetchConfig();
  }, []);

  if (!apiUrl) {
    return null;
  }

  return <Context.Provider value={apiUrl}>{children}</Context.Provider>;
}

export function useConsoleApiUrl(): string {
  const apiUrl = useContext(Context);
  if (!apiUrl) {
    throw new Error(
      "useConsoleApiUrl must be used within a ConsoleApiUrlProvider",
    );
  }
  return apiUrl;
}
