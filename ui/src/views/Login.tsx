import React, { Dispatch, FC, SetStateAction } from 'react'

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
import { useMutation } from '@connectrpc/connect-query'

import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { LoginViews } from '@/lib/views'
import useSettings from '@/lib/settings'

interface LoginProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const Login: FC<LoginProps> = ({ setView }) => {
  const settings = useSettings()

  const createIntermediateSessionMutation = useMutation(
    createIntermediateSession,
  )
  const googleOAuthRedirectUrlMutation = useMutation(getGoogleOAuthRedirectURL)
  const microsoftOAuthRedirectUrlMutation = useMutation(
    getMicrosoftOAuthRedirectURL,
  )

  const handleGoogleOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    try {
      // this sets a cookie that subsequent requests use
      await createIntermediateSessionMutation.mutateAsync({})
    } catch (error) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(error)
    }

    try {
      const { url } = await googleOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/google-oauth-callback`,
      })

      window.location.href = url
    } catch (error) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(error)
    }
  }

  const handleMicrosoftOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    try {
      // this sets a cookie that subsequent requests use
      await createIntermediateSessionMutation.mutateAsync({})
    } catch (error) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(error)
    }

    try {
      const { url } = await microsoftOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
      })

      window.location.href = url
    } catch (error) {
      // TODO: Handle errors on screen once an error handling strategy is in place.
      console.error(error)
    }
  }

  return (
    <>
      <Title title="Login" />

      <Card className="w-[clamp(320px,50%,420px)]">
        <CardHeader>
          {(settings?.logInWithGoogleEnabled ||
            settings?.logInWithMicrosoftEnabled) && (
            <CardTitle className="text-center uppercase text-foreground font-semibold text-sm tracking-wide mt-2">
              Log In with oAuth
            </CardTitle>
          )}
        </CardHeader>

        <CardContent className="flex flex-col items-center w-full">
          {settings?.logInWithGoogleEnabled && (
            <OAuthButton
              className="mb-4 w-[clamp(240px,50%,100%)]"
              method={OAuthMethods.google}
              onClick={handleGoogleOAuthLogin}
              variant="outline"
            />
          )}
          {settings?.logInWithMicrosoftEnabled && (
            <OAuthButton
              className="w-[clamp(240px,50%,100%)]"
              method={OAuthMethods.microsoft}
              onClick={handleMicrosoftOAuthLogin}
              variant="outline"
            />
          )}

          {(settings?.logInWithGoogleEnabled ||
            settings?.logInWithMicrosoftEnabled) && <TextDivider text="or" />}

          <EmailForm setView={setView} />
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </>
  )
}

export default Login
