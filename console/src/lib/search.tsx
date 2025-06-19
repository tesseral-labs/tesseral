import React, {
  Dispatch,
  SetStateAction,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

interface GlobalSearchContext {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const GlobalSearchContext = createContext<GlobalSearchContext | undefined>(
  undefined,
);

export const useGlobalSearch = () => {
  const context = useContext(GlobalSearchContext);
  if (!context) {
    throw new Error(
      "useGlobalSearch must be used within a GlobalSearchProvider",
    );
  }
  return context;
};

export function GlobalSearchProvider({ children }: React.PropsWithChildren) {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      const isCmdK =
        (event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k";
      if (isCmdK) {
        event.preventDefault();
        setOpen((prev) => !prev);
      }
    };

    window?.addEventListener("keydown", handleKeyDown);
    return () => window?.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <GlobalSearchContext.Provider value={{ open, setOpen }}>
      {children}
    </GlobalSearchContext.Provider>
  );
}
