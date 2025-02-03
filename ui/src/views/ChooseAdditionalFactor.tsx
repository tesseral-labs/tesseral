import React, { Dispatch, FC, SetStateAction, useEffect } from 'react'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

import { useIntermediateOrganization } from '@/lib/auth'
import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { LoginLayouts, LoginViews } from '@/lib/views'

interface ChooseAdditionalFactorProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const ChooseAdditionalFactor: FC<ChooseAdditionalFactorProps> = ({
  setView,
}) => {
  const org = useIntermediateOrganization()
  const layout = useLayout()

  const hasSecondFactor = org?.userHasPasskey || org?.userHasAuthenticatorApp

  return (
    <Card
      className={cn(
        'w-full max-w-sm',
        layout !== LoginLayouts.Centered && 'shadow-none border-0',
      )}
    >
      <CardHeader>
        <CardTitle className="text-center">Choose additional factor</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-2">
        {org?.logInWithPasskeyEnabled &&
          (!hasSecondFactor || org.userHasPasskey) && (
            <Button
              className="w-full"
              onClick={() =>
                setView(
                  org.userHasPasskey
                    ? LoginViews.VerifyPasskey
                    : LoginViews.RegisterPasskey,
                )
              }
            >
              Continue with Passkey
            </Button>
          )}
        {org?.logInWithAuthenticatorAppEnabled &&
          (!hasSecondFactor || org.userHasAuthenticatorApp) && (
            <Button
              className="w-full"
              onClick={() =>
                setView(
                  org.userHasAuthenticatorApp
                    ? LoginViews.VerifyAuthenticatorApp
                    : LoginViews.RegisterAuthenticatorApp,
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

export default ChooseAdditionalFactor
