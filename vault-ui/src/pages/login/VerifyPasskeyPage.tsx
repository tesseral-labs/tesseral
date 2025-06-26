import { useMutation } from "@connectrpc/connect-query";
import React, { useCallback, useEffect } from "react";

import { Title } from "@/components/core/Title";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  issuePasskeyChallenge,
  verifyPasskey,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";
import { base64urlEncode } from "@/lib/utils";

export function VerifyPasskeyPage() {
  const { mutateAsync: issuePasskeyChallengeAsync } = useMutation(
    issuePasskeyChallenge,
  );
  const { mutateAsync: verifyPasskeyAsync } = useMutation(verifyPasskey);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  const handleVerifyPasskey = useCallback(async () => {
    const challengeResponse = await issuePasskeyChallengeAsync({});

    const allowCredentials = challengeResponse.credentialIds.map(
      (id) =>
        ({
          id: new Uint8Array(id).buffer,
          type: "public-key",
          transports: ["hybrid", "internal", "nfc", "usb"],
        }) as PublicKeyCredentialDescriptor,
    );

    const requestOptions: PublicKeyCredentialRequestOptions = {
      challenge: new Uint8Array(challengeResponse.challenge).buffer,
      allowCredentials,
      rpId: challengeResponse.rpId,
      userVerification: "preferred",
      timeout: 60000,
    };
    const credential = (await navigator.credentials.get({
      publicKey: requestOptions,
    })) as PublicKeyCredential;

    const response = credential.response as AuthenticatorAssertionResponse;

    await verifyPasskeyAsync({
      authenticatorData: base64urlEncode(response.authenticatorData),
      clientDataJson: base64urlEncode(response.clientDataJSON),
      credentialId: new Uint8Array(credential.rawId),
      signature: base64urlEncode(response.signature),
    });

    redirectNextLoginFlowPage();
  }, [
    issuePasskeyChallengeAsync,
    redirectNextLoginFlowPage,
    verifyPasskeyAsync,
  ]);

  useEffect(() => {
    void handleVerifyPasskey();
  }, [handleVerifyPasskey]);

  return (
    <LoginFlowCard>
      <Title title="Verify passkey" />
      <CardHeader>
        <CardTitle>Verify passkey</CardTitle>
        <CardDescription>
          To continue logging in, follow the instructions on your authenticator.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button onClick={handleVerifyPasskey}>Verify passkey</Button>
      </CardContent>
    </LoginFlowCard>
  );
}
