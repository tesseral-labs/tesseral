import React, { FC } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useNavigate } from 'react-router'
import { useMutation } from '@connectrpc/connect-query'
import { useLayout } from '@/lib/settings'
import { cn } from '@/lib/utils'
import { LoginLayouts } from '@/lib/views'

const VerifyPasskey: FC = () => {
  const layout = useLayout()
  const navigate = useNavigate()

  return (
    <Card
      className={cn(
        'w-full max-w-sm',
        layout !== LoginLayouts.Centered && 'shadow-none border-0',
      )}
    >
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
