import { useMutation } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS } from "input-otp";
import QRCode from "qrcode";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Title } from "@/components/Title";
import { LoginFlowCard } from "@/components/login/LoginFlowCard";
import { Button } from "@/components/ui/button";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "@/components/ui/input-otp";
import {
  getAuthenticatorAppOptions,
  registerAuthenticatorApp,
} from "@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery";
import { useRedirectNextLoginFlowPage } from "@/hooks/use-redirect-next-login-flow-page";

const schema = z.object({
  totpCode: z.string().length(6),
});

export function RegisterAuthenticatorAppPage() {
  const { mutateAsync: getAuthenticatorAppOptionsAsync } = useMutation(
    getAuthenticatorAppOptions,
  );
  const [qrCode, setQRCode] = useState("");

  useEffect(() => {
    (async () => {
      const { otpauthUri } = await getAuthenticatorAppOptionsAsync({});
      setQRCode(
        await QRCode.toDataURL(otpauthUri, {
          errorCorrectionLevel: "high",
        }),
      );
    })();
  }, [getAuthenticatorAppOptionsAsync]);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      totpCode: "",
    },
  });

  const [recoveryCodes, setRecoveryCodes] = useState<string[] | undefined>();
  const { mutateAsync: registerAuthenticatorAppAsync } = useMutation(
    registerAuthenticatorApp,
  );

  async function handleSubmit(values: z.infer<typeof schema>) {
    const { recoveryCodes } = await registerAuthenticatorAppAsync({
      totpCode: values.totpCode,
    });
    setRecoveryCodes(recoveryCodes);
  }

  async function handleCopy() {
    await navigator.clipboard.writeText(recoveryCodes!.join("\n"));
    toast.success("Copied recovery codes to clipboard");
  }

  const redirectNextLoginFlowPage = useRedirectNextLoginFlowPage();

  async function handleFinish() {
    redirectNextLoginFlowPage();
  }

  return recoveryCodes ? (
    <LoginFlowCard>
      <Title title="Register authenticator app" />
      <CardHeader>
        <CardTitle>Authenticator app recovery codes</CardTitle>
        <CardDescription>
          Keep these recovery codes in a private place.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="max-w-full overflow-x-auto bg-muted rounded-md">
          <div className="p-2 font-mono text-xs">
            {recoveryCodes.map((recoveryCode, i) => (
              <div key={i}>{recoveryCode}</div>
            ))}
          </div>
        </div>

        <p className="mt-2 text-sm text-muted-foreground">
          Each code can only be used once to sign in if you lose access to your
          authenticator app.
        </p>

        <Button variant="outline" onClick={handleCopy} className="mt-4 w-full">
          Copy recovery codes
        </Button>
        <Button className="mt-2 w-full" onClick={handleFinish}>
          Finish logging in
        </Button>
      </CardContent>
    </LoginFlowCard>
  ) : (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Set up authenticator app</CardTitle>
        <CardDescription>
          Scan the QR Code below with your authenticator app to continue logging
          in.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <img src={qrCode} className="w-full" />

        <Form {...form}>
          <form className="mt-4" onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="totpCode"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>One-Time Password</FormLabel>
                  <FormControl>
                    <InputOTP
                      pattern={REGEXP_ONLY_DIGITS}
                      maxLength={6}
                      {...field}
                    >
                      <InputOTPGroup>
                        <InputOTPSlot index={0} />
                        <InputOTPSlot index={1} />
                        <InputOTPSlot index={2} />
                        <InputOTPSlot index={3} />
                        <InputOTPSlot index={4} />
                        <InputOTPSlot index={5} />
                      </InputOTPGroup>
                    </InputOTP>
                  </FormControl>
                  <FormDescription>
                    After you've scanned the QR Code, enter a six-digit code
                    from your authenticator app.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full">
              Set up authenticator app
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
