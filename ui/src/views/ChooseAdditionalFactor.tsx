import React, { Dispatch, FC, SetStateAction } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { LoginViews } from '@/lib/views'
import { Button } from '@/components/ui/button'
import { useIntermediateOrganization, useOrganization } from '@/lib/auth'
import { useQuery } from '@connectrpc/connect-query'

interface ChooseAdditionalFactorProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const ChooseAdditionalFactor: FC<ChooseAdditionalFactorProps> = ({
  setView,
}) => {
  const org = useIntermediateOrganization()

  return (
    <Card className="w-full max-w-sm">
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
