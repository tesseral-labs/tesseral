import { useMutation } from "@connectrpc/connect-query";
import React, { useCallback, useEffect } from "react";

import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  getPasskeyOptions,
  registerPasskey,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";
import { base64urlEncode } from "@/lib/utils";
import { Title } from "@/components/Title";

export function RegisterPasskeyPage() {
  const { mutateAsync: getPasskeyOptionsAsync } =
    useMutation(getPasskeyOptions);
  const { mutateAsync: registerPasskeyAsync } = useMutation(registerPasskey);
  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  const handleRegisterPasskey = useCallback(async () => {
    const passkeyOptions = await getPasskeyOptionsAsync({});
    const credentialOptions: PublicKeyCredentialCreationOptions = {
      challenge: new Uint8Array([0]).buffer,
      rp: {
        id: passkeyOptions.rpId,
        name: passkeyOptions.rpName,
      },
      user: {
        id: new TextEncoder().encode(passkeyOptions.userId).buffer,
        name: passkeyOptions.userDisplayName,
        displayName: passkeyOptions.userDisplayName,
      },
      pubKeyCredParams: [
        { type: "public-key", alg: -7 }, // ECDSA with SHA-256
        { type: "public-key", alg: -257 }, // RSA with SHA-256
      ],
      timeout: 60000,
      attestation: "direct",
    };

    const credential = (await navigator.credentials.create({
      publicKey: credentialOptions,
    })) as PublicKeyCredential;

    if (!credential) {
      throw new Error("No credential returned");
    }

    await registerPasskeyAsync({
      rpId: passkeyOptions.rpId,
      attestationObject: base64urlEncode(
        (credential.response as AuthenticatorAttestationResponse)
          .attestationObject,
      ),
    });

    redirectNextLoginFlowPage();
  }, [getPasskeyOptionsAsync, registerPasskeyAsync, redirectNextLoginFlowPage]);

  useEffect(() => {
    void handleRegisterPasskey();
  }, [handleRegisterPasskey]);

  return (
    <LoginFlowCard>
      <Title title="Register a passkey" />
      <CardHeader>
        <CardTitle>Register a passkey</CardTitle>
        <CardDescription>
          Follow the instructions on your authenticator to create a new passkey.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Button onClick={handleRegisterPasskey}>Register passkey</Button>
      </CardContent>
    </LoginFlowCard>
  );
}
