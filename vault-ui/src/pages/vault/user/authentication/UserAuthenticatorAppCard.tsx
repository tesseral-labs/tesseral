import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { REGEXP_ONLY_DIGITS } from "input-otp";
import { Plus } from "lucide-react";
import QRCode from "qrcode";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { Label } from "@/components/ui/label";
import {
  getAuthenticatorAppOptions,
  registerAuthenticatorApp,
  whoami,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function UserAuthenticatorAppCard() {
  const { data: whoamiResponse } = useQuery(whoami);

  const user = whoamiResponse?.user;

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Authenticator App</CardTitle>
        <CardDescription>
          Authenticator Apps allow you to log in using a one-time code from an
          app on your device.
        </CardDescription>
        <CardAction>
          <RegisterAuthenticatorAppButton />
        </CardAction>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label>Status</Label>
          {user?.hasAuthenticatorApp ? (
            <Badge>Registered</Badge>
          ) : (
            <Badge variant="secondary">Not Registered</Badge>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  totpCode: z.string().length(6),
});

export function RegisterAuthenticatorAppButton() {
  const { data: whoamiResponse, refetch } = useQuery(whoami);
  const getAuthenticatorAppOptionsMutation = useMutation(
    getAuthenticatorAppOptions,
  );
  const [qrCode, setQRCode] = useState("");
  const [registerOpen, setRegisterOpen] = useState(false);
  const [recoveryOpen, setRecoveryOpen] = useState(false);
  const [recoveryCodes, setRecoveryCodes] = useState<string[] | undefined>();
  const registerAuthenticatorAppMutation = useMutation(
    registerAuthenticatorApp,
  );

  const user = whoamiResponse?.user;

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      totpCode: "",
    },
  });

  async function handleClick() {
    const { otpauthUri } = await getAuthenticatorAppOptionsMutation.mutateAsync(
      {},
    );
    setQRCode(
      await QRCode.toDataURL(otpauthUri, {
        errorCorrectionLevel: "high",
      }),
    );

    setRegisterOpen(true);
  }

  async function handleSubmit(values: z.infer<typeof schema>) {
    try {
      const { recoveryCodes } =
        await registerAuthenticatorAppMutation.mutateAsync({
          totpCode: values.totpCode,
        });
      setRecoveryCodes(recoveryCodes);
    } catch {
      toast.error("Failed to register authenticator app. Please try again.");
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
      <Dialog open={registerOpen} onOpenChange={setRegisterOpen}>
        <Button
          variant={user?.hasAuthenticatorApp ? "outline" : "default"}
          onClick={handleClick}
        >
          <Plus />
          {user?.hasAuthenticatorApp
            ? "Re-register authenticator app"
            : "Register authenticator app"}
        </Button>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Set up authenticator app</DialogTitle>
            <DialogDescription>
              Scan the QR Code below with your authenticator app.
            </DialogDescription>
          </DialogHeader>

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
                    <FormMessage />
                    <FormControl>
                      <InputOTP
                        pattern={REGEXP_ONLY_DIGITS}
                        maxLength={6}
                        {...field}
                      >
                        <InputOTPGroup className="w-full justify-center">
                          <InputOTPSlot index={0} />
                          <InputOTPSlot index={1} />
                          <InputOTPSlot index={2} />
                          <InputOTPSlot index={3} />
                          <InputOTPSlot index={4} />
                          <InputOTPSlot index={5} />
                        </InputOTPGroup>
                      </InputOTP>
                    </FormControl>
                  </FormItem>
                )}
              />

              <DialogFooter className="mt-8">
                <Button
                  variant="outline"
                  onClick={() => setRegisterOpen(false)}
                >
                  Cancel
                </Button>
                <Button type="submit">Set up authenticator app</Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

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
