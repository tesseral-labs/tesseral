import { useQuery } from '@tanstack/react-query'
import { DOGFOOD_PROJECT_ID } from '@/config'

export function useAccessToken(): string {
  const { data: accessToken } = useQuery({
    queryKey: ["refresh"],
    queryFn: async () => {
      const response = await fetch("http://api.tesseral.example.com/frontend/v1/access-token", {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-TODO-OpenAuth-Project-ID": DOGFOOD_PROJECT_ID,
        },
        body: "{}"
      })
      return (await response.json()).accessToken
    }
  })

  return accessToken
}
