import React, { Dispatch, FC, SetStateAction, useState } from 'react'

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
import TextDivider from '@/components/ui/test-divider'
import { useMutation } from '@connectrpc/connect-query'

import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { LoginViews } from '@/lib/views'
import useSettings from '@/lib/settings'
import { cn } from '@/lib/utils'

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

      <Card>
        <CardHeader>
          {(settings?.logInWithGoogleEnabled ||
            settings?.logInWithMicrosoftEnabled) && (
            <CardTitle className="text-center">Continue with OAuth</CardTitle>
          )}
        </CardHeader>

        <CardContent className="flex flex-col items-center w-full">
          <div
            className={cn(
              'grid gap-6',
              settings?.logInWithGoogleEnabled &&
                settings?.logInWithMicrosoftEnabled
                ? 'grid-cols-2'
                : 'grid-cols-1',
            )}
          >
            {settings?.logInWithGoogleEnabled && (
              <OAuthButton
                method={OAuthMethods.google}
                onClick={handleGoogleOAuthLogin}
                variant="outline"
              />
            )}
            {settings?.logInWithMicrosoftEnabled && (
              <OAuthButton
                method={OAuthMethods.microsoft}
                onClick={handleMicrosoftOAuthLogin}
                variant="outline"
              />
            )}
          </div>

          {(settings?.logInWithGoogleEnabled ||
            settings?.logInWithMicrosoftEnabled) && (
            <TextDivider variant="wide">or continue with email</TextDivider>
          )}

          <EmailForm setView={setView} />
        </CardContent>
      </Card>
    </>
  )
}

export default Login
