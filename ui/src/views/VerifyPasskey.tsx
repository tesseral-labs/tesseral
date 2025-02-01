import React, { FC } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'

const VerifyPasskey: FC = () => {
  const navigate = useNavigate()

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle className="text-center">Verify Passkey</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-center text-sm text-muted-foreground">
          Follow the prompts on your device to continue logging in with your
          Passkey.
        </p>
      </CardContent>
    </Card>
  )
}

export default VerifyPasskey
