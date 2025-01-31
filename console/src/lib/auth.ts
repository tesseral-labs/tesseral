import { useMutation } from '@tanstack/react-query'
import { base64Decode } from './utils'
import { useNavigate } from 'react-router'
import { createContext, useContext, useState } from 'react'
import { useQuery } from '@connectrpc/connect-query'
import { whoami } from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { IntermediateSession } from '@/gen/openauth/intermediate/v1/intermediate_pb'

export const useIntermediateSession = (): IntermediateSession | undefined => {
  const { data: whoAmiRes } = useQuery(whoami)

  return whoAmiRes?.intermediateSession
}

// how far in advance of its expiration an access token gets refreshed
const ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS = 10

interface SessionAccessTokenClaims {
  exp: number
  iat: number
  nbf: number
  organization: SessionOrganizationClaims
  project: SessionProjectClaims
  user: SessionUserClaims
}

interface SessionOrganizationClaims {
  createTime: string
  displayName: string
  id: string
  logInWithGoogleEnabled: boolean
  logInWithMicrosoftEnabled: boolean
  logInWithPasswordEnabled: boolean
  overrideLogInMethods: boolean
  samlEnabled: boolean
  updateTime: string
}

interface SessionProjectClaims {
  authDomain: string
  createTime: string
  displayName: string
  id: string
  logInWithGoogleEnabled: boolean
  logInWithMicrosoftEnabled: boolean
  logInWithPasswordEnabled: boolean
  updateTime: string
}

interface SessionUserClaims {
  createTime: string
  email: string
  googleUserId?: string
  id: string
  microsoftUserId?: string
  owner: boolean
  updateTime: string
}

export const useSession = (): SessionAccessTokenClaims | undefined => {
  const navigate = useNavigate()

  // read the access token from local storage
  const [accessToken, setAccessToken] = useState<string | null>(
    localStorage.getItem(`accessToken`),
  )

  // mutation for refreshing the user's access token if necessary
  const refresh = useMutation({
    mutationKey: ['refresh'],
    mutationFn: async () => {
      const response = await fetch('/api/frontend/v1/refresh', {
        credentials: 'include',
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: '{}',
      })

      return (await response.json()).accessToken
    },
  })

  // if the user is fetched successfully, attempt to refresh the access token
  if (!accessToken || shouldRefresh(accessToken)) {
    if (!refresh.isPending) {
      refresh.mutate(undefined, {
        onSuccess: (accessToken) => {
          if (accessToken) {
            localStorage.setItem(`access_token`, accessToken)
            setAccessToken(accessToken)
          }
        },
      })
    }
  }

  if (refresh.isError) {
    console.error(refresh.error)
    navigate('/login')
    return
  }

  if (accessToken) {
    // Check if the access token is expired
    if (isAccessTokenExpired(accessToken) && !refresh.isPending) {
      navigate('/login')
      return
    }

    // parse the access token and return the user claims
    const claims = JSON.parse(
      base64Decode(accessToken.split('.')[1]),
    ) as SessionAccessTokenClaims

    return claims
  }

  return
}

const organizationContext = createContext<
  SessionOrganizationClaims | undefined
>(undefined)
export const OrganizationContextProvider = organizationContext.Provider
export const useOrganization = (): SessionOrganizationClaims | undefined => {
  return useContext(organizationContext)
}

const projectContext = createContext<SessionProjectClaims | undefined>(
  undefined,
)
export const ProjectContextProvider = projectContext.Provider
export const useProject = (): SessionProjectClaims | undefined => {
  return useContext(projectContext)
}

const userContext = createContext<SessionUserClaims | undefined>(undefined)
export const UserContextProvider = userContext.Provider
export const useUser = (): SessionUserClaims | undefined => {
  return useContext(userContext)
}

const isAccessTokenExpired = (accessToken: string): boolean => {
  return (
    parseAccessTokenExpiration(accessToken) <
    Math.floor(new Date().getTime() / 1000)
  )
}

const parseAccessTokenExpiration = (accessToken: string): number => {
  const claims = JSON.parse(
    base64Decode(accessToken.split('.')[1]),
  ) as SessionAccessTokenClaims
  return claims.exp
}

const shouldRefresh = (accessToken: string): boolean => {
  const refreshAt =
    parseAccessTokenExpiration(accessToken) -
    ACCESS_TOKEN_REFRESH_THRESHOLD_SECONDS
  const now = Math.floor(new Date().getTime() / 1000)
  return refreshAt < now
}
