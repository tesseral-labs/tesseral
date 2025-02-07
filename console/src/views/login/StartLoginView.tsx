import React, { Dispatch, FC } from 'react'
import { LoginView } from '@/lib/views'
import OAuthButton, { OAuthMethods } from '@/components/login/OAuthButton'
import { useMutation } from '@connectrpc/connect-query'
import {
  createIntermediateSession,
  getGoogleOAuthRedirectURL,
  getMicrosoftOAuthRedirectURL,
} from '@/gen/openauth/intermediate/v1/intermediate-IntermediateService_connectquery'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import TextDivider from '@/components/ui/text-divider'
import EmailForm from '@/components/login/EmailForm'
import { AuthType, useAuthType } from '@/lib/auth'
import { Title } from '@/components/Title'
import { parseErrorMessage } from '@/lib/errors'
import { toast } from 'sonner'
import useSettings from '@/lib/settings'

interface StartLoginViewProps {
  setView: Dispatch<React.SetStateAction<LoginView>>
}

const StartLoginView: FC<StartLoginViewProps> = ({ setView }) => {
  const authType = useAuthType()
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
      toast.error('Could not initiate log in', {
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
      toast.error('Could not log in with Google', {
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
      toast.error('Could not initiate log in', {
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
      toast.error('Could not log in with Microsoft', {
        description: message,
      })
    }
  }

  return (
    <>
      <Title title={authType === AuthType.SignUp ? 'Sign up' : 'Log in'} />

      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-center">
            {authType === AuthType.SignUp ? 'Sign up' : 'Log in'} with
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center w-full">
          <div className="w-full grid grid-cols-2 gap-6">
            {settings?.logInWithGoogle && (
              <OAuthButton
                className="w-full"
                method={OAuthMethods.google}
                onClick={handleGoogleOAuthLogin}
                variant="outline"
              />
            )}
            {settings?.logInWithMicrosoft && (
              <OAuthButton
                className="mt-4w-full"
                method={OAuthMethods.microsoft}
                onClick={handleMicrosoftOAuthLogin}
                variant="outline"
              />
            )}
          </div>

          {(settings?.logInWithEmail || settings?.logInWithSaml) && (
            <>
              <TextDivider>Or continue with email</TextDivider>

              <EmailForm
                disableLogInWithEmail={!settings?.logInWithEmail}
                skipListSAMLOrganizations={!settings?.logInWithSaml}
                setView={setView}
              />
            </>
          )}
        </CardContent>
      </Card>
    </>
  )
}

export default StartLoginView
