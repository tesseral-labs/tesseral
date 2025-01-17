import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'

interface User {
  id: string
  email: string
}

export function useUser(): User | undefined {
  const accessToken = useAccessToken()
  if (!accessToken) {
    return
  }

  const claims = JSON.parse(base64Decode(accessToken.split('.')[1]))
  return {
    id: claims.user.id,
    email: claims.user.email,
  }
}

export function useAccessToken(): string | undefined {
  const [accessToken, setAccessToken] = useLocalStorage('accessToken')
  const refresh = useMutation({
    mutationKey: ['refresh'],
    mutationFn: async () => {
      const response = await fetch(
        `https://auth.app.tesseral.example.com/api/frontend/v1/access-token`,
        {
          credentials: 'include',
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: '{}',
        },
      )
      return (await response.json()).accessToken
    },
  })

  if (!accessToken || shouldRefresh(accessToken)) {
    if (!refresh.isPending) {
      refresh.mutate(undefined, {
        onSuccess: (accessToken) => {
          setAccessToken(accessToken)
        },
      })
    }
  }

  return accessToken ?? undefined
}

// how far in advance of its expiration an access token gets refreshed
const ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS = 10

function shouldRefresh(accessToken: string): boolean {
  const refreshAt =
    parseAccessTokenExpiration(accessToken) -
    ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS
  const now = Math.floor(new Date().getTime() / 1000)
  return refreshAt < now
}

function parseAccessTokenExpiration(accessToken: string): number {
  return JSON.parse(base64Decode(accessToken.split('.')[1])).exp
}

function base64Decode(s: string): string {
  const binaryString = atob(s)

  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }

  return new TextDecoder().decode(bytes)
}

function useLocalStorage(
  key: string,
): [string | null, (_: string | null) => void] {
  const [value, setValue] = useState<string | null>(localStorage.getItem(key))

  return [
    value,
    (value) => {
      if (value === null) {
        localStorage.removeItem(key)
        setValue(null)
        return
      }

      localStorage.setItem(key, value)
      setValue(value)
    },
  ]
}
