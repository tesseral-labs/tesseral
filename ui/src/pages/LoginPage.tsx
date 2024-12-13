import React, { useEffect } from 'react'

import EmailForm from '@/components/EmailForm'
import OAuthButton, { OAuthMethods } from '@/components/OAuthButton'
import { Title } from '@/components/Title'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import TextDivider from '@/components/ui/TextDivider'
import { useQuery } from '@connectrpc/connect-query'

import { setIntermediateSessionToken } from '@/auth'
import {
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'

const LoginPage = () => {
  const googleOAuthRedirectUrlQuery = useQuery(getGoogleOAuthRedirectURL)
  const microsoftOAuthRedirectUrlQuery = useQuery(getMicrosoftOAuthRedirectURL)

  const [googleOAuthRedirectUrl, setGoogleOAuthRedirectUrl] = React.useState('')
  const [microsoftOAuthRedirectUrl, setMicrosoftOAuthRedirectUrl] =
    React.useState('')

  const handleGoogleOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    if (googleOAuthRedirectUrl) {
      window.location.href = googleOAuthRedirectUrl
      return
    }

    const response = await googleOAuthRedirectUrlQuery.refetch()

    if (response.isError) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(response.error)
      return
    }

    if (response.data) {
      setIntermediateSessionToken(response.data.intermediateSessionToken)
      window.location.href = response.data.url
    }
  }

  const handleMicrosoftOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    if (microsoftOAuthRedirectUrl) {
      window.location.href = microsoftOAuthRedirectUrl
      return
    }

    const response = await microsoftOAuthRedirectUrlQuery.refetch()

    if (response.isError) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(response.error)
      return
    }

    if (response.data) {
      setIntermediateSessionToken(response.data.intermediateSessionToken)
      window.location.href = response.data.url
    }
  }

  useEffect(() => {
    ;(async () => {
      if (googleOAuthRedirectUrlQuery.isError) {
        // TODO: Handle errors on screen once an error handling strategy is in place.
        console.error(
          `Error fetching Google OAuth redirect URL: ${googleOAuthRedirectUrlQuery.error}`,
        )
      }

      if (googleOAuthRedirectUrlQuery.data) {
        setIntermediateSessionToken(
          googleOAuthRedirectUrlQuery.data.intermediateSessionToken,
        )
        setGoogleOAuthRedirectUrl(googleOAuthRedirectUrlQuery.data.url)
      }
    })()
  }, [googleOAuthRedirectUrlQuery])

  useEffect(() => {
    ;(async () => {
      if (microsoftOAuthRedirectUrlQuery.isError) {
        // TODO: Handle errors on screen once an error handling strategy is in place.
        console.error(
          `Error fetching Microsoft OAuth redirect URL: ${microsoftOAuthRedirectUrlQuery.error}`,
        )
      }

      if (microsoftOAuthRedirectUrlQuery.data) {
        setIntermediateSessionToken(
          microsoftOAuthRedirectUrlQuery.data.intermediateSessionToken,
        )
        setMicrosoftOAuthRedirectUrl(microsoftOAuthRedirectUrlQuery.data.url)
      }
    })()
  }, [microsoftOAuthRedirectUrlQuery])

  return (
    <>
      <Title title="Login" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
            Log In with oAuth
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center w-full">
          <OAuthButton
            className="mb-4 w-[clamp(240px,50%,100%)]"
            method={OAuthMethods.google}
            onClick={handleGoogleOAuthLogin}
            variant="outline"
          />
          <OAuthButton
            className="w-[clamp(240px,50%,100%)]"
            method={OAuthMethods.microsoft}
            onClick={handleMicrosoftOAuthLogin}
            variant="outline"
          />

          <TextDivider text="or" />

          <EmailForm />
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default LoginPage
