import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useIntermediateOrganization } from '@/lib/auth'
import { LoginView } from '@/lib/views'
import React, { Dispatch, FC, SetStateAction } from 'react'

interface ChooseAdditionalFactorViewProps {
  setView: Dispatch<SetStateAction<LoginView>>
}

const ChooseAdditionalFactorView: FC<ChooseAdditionalFactorViewProps> = ({
  setView,
}) => {
  const org = useIntermediateOrganization()

  const hasSecondFactor = org?.userHasPasskey || org?.userHasAuthenticatorApp

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle className="text-center">Choose additional factor</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-2">
        {org?.logInWithPasskey && (!hasSecondFactor || org.userHasPasskey) && (
          <Button
            className="w-full"
            onClick={() =>
              setView(
                org.userHasPasskey
                  ? LoginView.VerifyPasskey
                  : LoginView.RegisterPasskey,
              )
            }
          >
            Continue with Passkey
          </Button>
        )}
        {org?.logInWithAuthenticatorApp &&
          (!hasSecondFactor || org.userHasAuthenticatorApp) && (
            <Button
              className="w-full"
              onClick={() =>
                setView(
                  org.userHasAuthenticatorApp
                    ? LoginView.VerifyAuthenticatorApp
                    : LoginView.RegisterAuthenticatorApp,
                )
              }
            >
              Continue with Authenticator app
            </Button>
          )}
      </CardContent>
    </Card>
  )
}

export default ChooseAdditionalFactorView
