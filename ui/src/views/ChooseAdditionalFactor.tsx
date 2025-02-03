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
        {org?.logInWithPasskeyEnabled && (
          <Button
            className="w-full"
            onClick={() => setView(LoginViews.RegisterPasskey)}
          >
            Continue with Passkey
          </Button>
        )}
        {org?.logInWithAuthenticatorAppEnabled && (
          <Button
            className="w-full"
            onClick={() => setView(LoginViews.RegisterAuthenticatorApp)}
          >
            Continue with Authenticator app
          </Button>
        )}
      </CardContent>
    </Card>
  )
}

export default ChooseAdditionalFactor
