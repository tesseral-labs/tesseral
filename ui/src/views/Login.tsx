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
import { LoginLayouts, LoginViews } from '@/lib/views'
import useSettings, { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'

interface LoginProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const Login: FC<LoginProps> = ({ setView }) => {
  const layout = useLayout()
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
      const message = parseErrorMessage(error)

      toast.error('Could not initialize new session', {
        description: message,
      })
    }

    try {
      const { url } = await googleOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/google-oauth-callback`,
      })

      window.location.href = url
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error('Could not get Google OAuth redirect URL', {
        description: message,
      })
    }
  }

  const handleMicrosoftOAuthLogin = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()

    try {
      // this sets a cookie that subsequent requests use
      await createIntermediateSessionMutation.mutateAsync({})
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error('Could not initialize new session', {
        description: message,
      })
    }

    try {
      const { url } = await microsoftOAuthRedirectUrlMutation.mutateAsync({
        redirectUrl: `${window.location.origin}/microsoft-oauth-callback`,
      })

      window.location.href = url
    } catch (error) {
      const message = parseErrorMessage(error)

      toast.error('Could not get Microsoft OAuth redirect URL', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title="Login" />

      <Card
        className={cn(
          'w-full max-w-sm',
          layout !== LoginLayouts.Centered && 'shadow-none border-0',
        )}
      >
        <CardHeader>
          {(settings?.logInWithGoogle || settings?.logInWithMicrosoft) && (
            <CardTitle className="text-center">Log in with</CardTitle>
          )}
        </CardHeader>

        <CardContent className="flex flex-col items-center w-full">
          <div
            className={cn(
              'w-full grid gap-6',
              settings?.logInWithGoogle && settings?.logInWithMicrosoft
                ? 'grid-cols-2'
                : 'grid-cols-1',
            )}
          >
            {settings?.logInWithGoogle && (
              <OAuthButton
                method={OAuthMethods.google}
                onClick={handleGoogleOAuthLogin}
                variant="outline"
              />
            )}
            {settings?.logInWithMicrosoft && (
              <OAuthButton
                method={OAuthMethods.microsoft}
                onClick={handleMicrosoftOAuthLogin}
                variant="outline"
              />
            )}
          </div>

          {(settings?.logInWithGoogle || settings?.logInWithMicrosoft) && (
            <TextDivider
              variant={layout !== LoginLayouts.Centered ? 'wider' : 'wide'}
            >
              or continue with email
            </TextDivider>
          )}

          {(settings?.logInWithEmail || settings?.logInWithSaml) && (
            <EmailForm
              skipListSAMLOrganizations={!settings?.logInWithSaml}
              setView={setView}
            />
          )}
        </CardContent>
      </Card>
    </>
  )
}

export default Login
