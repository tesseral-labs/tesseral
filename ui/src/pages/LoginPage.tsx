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
import React from 'react'

const LoginPage = () => {
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
            variant="outline"
          />
          <OAuthButton
            className="w-[clamp(240px,50%,100%)]"
            method={OAuthMethods.microsoft}
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
