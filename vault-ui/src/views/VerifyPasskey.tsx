import React, { FC, useEffect } from 'react';
import { useNavigate } from 'react-router';
import { useMutation } from '@connectrpc/connect-query';
import { toast } from 'sonner';

import { useLayout } from '@/lib/settings';
import { base64urlEncode, cn } from '@/lib/utils';
import { LoginLayouts } from '@/lib/views';
import { parseErrorMessage } from '@/lib/errors';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  issuePasskeyChallenge,
  verifyPasskey,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { setAccessToken, setRefreshToken } from '@/auth';
import {
  useIntermediateExchangeAndRedirect
} from '@/hooks/use-intermediate-exchange-and-redirect';

const VerifyPasskey: FC = () => {
  const layout = useLayout();
  const navigate = useNavigate();

  const issuePasskeyChallengeMutation = useMutation(issuePasskeyChallenge);
  const verifyPasskeyMutation = useMutation(verifyPasskey);
  const intermediateExchangeAndRedirect = useIntermediateExchangeAndRedirect()

  const authenticateWithPasskey = async () => {
    try {
      const challengeResponse = await issuePasskeyChallengeMutation.mutateAsync(
        {},
      );

      const allowCredentials = challengeResponse.credentialIds.map(
        (id) =>
          ({
            id: new Uint8Array(id).buffer,
            type: 'public-key',
            transports: ['hybrid', 'internal', 'nfc', 'usb'],
          }) as PublicKeyCredentialDescriptor,
      );

      const requestOptions: PublicKeyCredentialRequestOptions = {
        challenge: new Uint8Array(challengeResponse.challenge).buffer,
        allowCredentials,
        rpId: challengeResponse.rpId,
        userVerification: 'preferred',
        timeout: 60000,
      };
      const credential = (await navigator.credentials.get({
        publicKey: requestOptions,
      })) as PublicKeyCredential;

      const response = credential.response as AuthenticatorAssertionResponse;

      await verifyPasskeyMutation.mutateAsync({
        authenticatorData: base64urlEncode(response.authenticatorData),
        clientDataJson: base64urlEncode(response.clientDataJSON),
        credentialId: new Uint8Array(credential.rawId),
        signature: base64urlEncode(response.signature),
      });

      intermediateExchangeAndRedirect();
    } catch (error) {
      const message = parseErrorMessage(error);
      toast.error('Could not verify passkey', {
        description: message,
      });
    }
  };

  useEffect(() => {
    void (async () => {
      await authenticateWithPasskey();
    })();
  }, []);

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
  );
};

export default VerifyPasskey;
