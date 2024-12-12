import React, {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from 'react'

export type Config = {
  API_URL: string
}

const Context = createContext<Config>(undefined as any)

export const ConfigProvider = ({ children }: { children?: ReactNode }) => {
  const [config, setConfig] = useState<Config | undefined>()

  useEffect(() => {
    ;(async () => {
      const res = await (await fetch('/config.json')).json()
      setConfig(res)
    })()
  }, [])

  if (!config) {
    return
  }

  return <Context.Provider value={config}>{children}</Context.Provider>
}

export const useConfig = (): Config => {
  return useContext(Context)
}
