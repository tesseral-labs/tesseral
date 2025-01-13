import { useQuery } from '@tanstack/react-query'
import { API_URL, DOGFOOD_PROJECT_ID } from '@/config'

export function useAccessToken(): string {
  const { data: accessToken } = useQuery({
    queryKey: ['refresh'],
    queryFn: async () => {
      const response = await fetch(`${API_URL}/api/frontend/v1/access-token`, {
        credentials: 'include',
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-TODO-OpenAuth-Project-ID': DOGFOOD_PROJECT_ID,
        },
        body: '{}',
      })
      return (await response.json()).accessToken
    },
  })

  return accessToken
}
