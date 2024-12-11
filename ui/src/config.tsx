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

const configData: Config = {
  API_URL: '__API_URL__',
}

const Context = createContext<Config>(undefined as any)

export const ConfigProvider = ({ children }: { children?: ReactNode }) => {
  const [config, setConfig] = useState<Config | undefined>()

  useEffect(() => {
    setConfig(configData)
  }, [])

  if (!config) {
    return
  }

  return <Context.Provider value={config}>{children}</Context.Provider>
}

export const useConfig = (): Config => {
  return useContext(Context)
}
