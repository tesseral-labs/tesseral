import React, { Dispatch, FC } from 'react'
import { LoginView } from '@/lib/views'
import OAuthButton, { OAuthMethods } from '@/components/login/OAuthButton'
import { useMutation } from '@connectrpc/connect-query'
import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import TextDivider from '@/components/ui/text-divider'
import EmailForm from '@/components/login/EmailForm'

interface StartLoginViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const StartLoginView: FC<StartLoginViewProps> = ({ setView }) => {
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
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle className="text-center">Log in with</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col items-center w-full">
        <div className="w-full grid grid-cols-2 gap-6">
          <OAuthButton
            className="w-full"
            method={OAuthMethods.google}
            onClick={handleGoogleOAuthLogin}
            variant="outline"
          />
          <OAuthButton
            className="mt-4w-full"
            method={OAuthMethods.microsoft}
            onClick={handleMicrosoftOAuthLogin}
            variant="outline"
          />
        </div>

        <TextDivider>Or continue with email</TextDivider>

        <EmailForm setView={setView} />
      </CardContent>
    </Card>
  )
}

export default StartLoginView
