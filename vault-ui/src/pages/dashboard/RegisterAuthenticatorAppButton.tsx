import { ConnectError } from "@connectrpc/connect";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS } from "input-otp";
import QRCode from "qrcode";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
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
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

const schema = z.object({
  totpCode: z.string().length(6),
});

export function RegisterAuthenticatorAppButton() {
  const { data: whoamiResponse, refetch } = useQuery(whoami);
  const { mutateAsync: getAuthenticatorAppOptionsAsync } = useMutation(
    getAuthenticatorAppOptions,
  );
  const [qrCode, setQRCode] = useState("");
  const [registerOpen, setRegisterOpen] = useState(false);
  const [recoveryOpen, setRecoveryOpen] = useState(false);

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      totpCode: "",
    },
  });

  async function handleClick() {
    const { otpauthUri } = await getAuthenticatorAppOptionsAsync({});
    setQRCode(
      await QRCode.toDataURL(otpauthUri, {
        errorCorrectionLevel: "high",
      }),
    );

    setRegisterOpen(true);
  }

  const [recoveryCodes, setRecoveryCodes] = useState<string[] | undefined>();
  const { mutateAsync: registerAuthenticatorAppAsync } = useMutation(
    registerAuthenticatorApp,
  );

  async function handleSubmit(values: z.infer<typeof schema>) {
    try {
      const { recoveryCodes } = await registerAuthenticatorAppAsync({
        totpCode: values.totpCode,
      });
      setRecoveryCodes(recoveryCodes);
    } catch (e) {
      if (
        e instanceof ConnectError &&
        e.message === "[invalid_argument] incorrect_totp_code"
      ) {
        form.setError("totpCode", {
          message: "Incorrect code.",
        });
        return;
      }

      throw e; // unhandled
    }

    await refetch();
    setRegisterOpen(false);
    setRecoveryOpen(true);
  }

  async function handleCopy() {
    await navigator.clipboard.writeText(recoveryCodes!.join("\n"));
    toast.success("Copied recovery codes to clipboard");
  }

  function handleDone() {
    toast.success("Authenticator app registered");
  }

  return (
    <>
      <AlertDialog open={registerOpen} onOpenChange={setRegisterOpen}>
        <Button variant="outline" onClick={handleClick}>
          {whoamiResponse?.user?.hasAuthenticatorApp
            ? "Re-register authenticator app"
            : "Register authenticator app"}
        </Button>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Set up authenticator app</AlertDialogTitle>
            <AlertDialogDescription>
              Scan the QR Code below with your authenticator app.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <img src={qrCode} className="w-full" />

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
              <FormField
                control={form.control}
                name="totpCode"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>One-Time Password</FormLabel>
                    <FormDescription>
                      After you've scanned the QR Code, enter a six-digit code
                      from your authenticator app.
                    </FormDescription>
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

                    <FormMessage />
                  </FormItem>
                )}
              />

              <AlertDialogFooter className="mt-8">
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Set up authenticator app</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={recoveryOpen} onOpenChange={setRecoveryOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              Authenticator app recovery codes
            </AlertDialogTitle>
            <AlertDialogDescription>
              Keep these recovery codes in a private place.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="flex justify-center">
            <div className="p-2 bg-muted rounded-md font-mono text-xs">
              {recoveryCodes?.map((recoveryCode, i) => (
                <div key={i}>{recoveryCode}</div>
              ))}
            </div>
          </div>

          <p className="mt-2 text-sm text-muted-foreground">
            Each code can only be used once to sign in if you lose access to
            your authenticator app.
          </p>

          <Button
            variant="outline"
            onClick={handleCopy}
            className="mt-4 w-full"
          >
            Copy recovery codes
          </Button>

          <AlertDialogFooter>
            <AlertDialogAction onClick={handleDone}>Done</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
