import React, { Dispatch, FC, SetStateAction } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { LoginViews } from '@/lib/views'
import { Button } from '@/components/ui/button'

interface ChooseAdditionalFactorProps {
  setView: Dispatch<SetStateAction<LoginViews>>
}

const ChooseAdditionalFactor: FC<ChooseAdditionalFactorProps> = ({
  setView,
}) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Choose additional factor</CardTitle>
        <CardContent>
          <Button
            className="block w-full p-4 border-b"
            onClick={() => setView(LoginViews.RegisterPasskey)}
            variant="outline"
          >
            Continue with Passkey
          </Button>
          <Button
            className="block w-full p-4 border-b"
            onClick={() => setView(LoginViews.RegisterAuthenticatorApp)}
            variant="outline"
          >
            Continue with Authenticator app
          </Button>
        </CardContent>
      </CardHeader>
    </Card>
  )
}

export default ChooseAdditionalFactor
