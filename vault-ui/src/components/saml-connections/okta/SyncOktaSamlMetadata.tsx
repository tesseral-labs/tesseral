import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  getSAMLConnection,
  updateSAMLConnection,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

interface OktaMetadata {
  idpEntityId: string;
  idpRedirectUrl: string;
  idpX509Certificate: string;
}

const schema = z.object({
  oktaMetadataUrl: z.string().url("Must be a valid URL"),
});

export function SyncOktaSamlMetadata() {
  const parser = new DOMParser();
  const { samlConnectionId } = useParams();
  const navigate = useNavigate();

  const { refetch } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });
  const updateSamlConnectionMutation = useMutation(updateSAMLConnection);

  const [oktaMetadata, setOktaMetadata] = useState<OktaMetadata>();

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      oktaMetadataUrl: "",
    },
  });

  async function handleSubmit() {
    try {
      await updateSamlConnectionMutation.mutateAsync({
        id: samlConnectionId,
        samlConnection: {
          idpEntityId: oktaMetadata?.idpEntityId,
          idpRedirectUrl: oktaMetadata?.idpRedirectUrl,
          idpX509Certificate: oktaMetadata?.idpX509Certificate,
        },
      });
      await refetch();
      form.reset();
      toast.success("Okta metadata synced successfully.");
      navigate(
        `/organization/saml-connections/${samlConnectionId}/setup/okta/users`,
      );
    } catch {
      toast.error("Failed to sync Okta metadata. Please try again.");
      return;
    }
  }

  async function handleUrlChange(e: React.ChangeEvent<HTMLInputElement>) {
    const url = e.target.value;
    form.setValue("oktaMetadataUrl", url || "", {
      shouldValidate: true,
      shouldDirty: true,
      shouldTouch: true,
    });

    if (!url || url.trim() === "") {
      return;
    }

    const response = await fetch(url);
    if (!response.ok) {
      form.setError("oktaMetadataUrl", {
        type: "manual",
        message: "Failed to fetch metadata from the provided URL.",
      });
      return;
    }

    const xmlText = await response.text();
    const xml = parser.parseFromString(xmlText, "application/xml");

    const mdNS = "urn:oasis:names:tc:SAML:2.0:metadata";
    const dsNS = "http://www.w3.org/2000/09/xmldsig#";

    const entityID = xml.documentElement.getAttribute("entityID") ?? undefined;

    const ssoServices = xml.getElementsByTagNameNS(mdNS, "SingleSignOnService");
    let ssoUrl: string | undefined;
    for (let i = 0; i < ssoServices.length; i++) {
      const el = ssoServices[i];
      if (
        el.getAttribute("Binding") ===
        "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
      ) {
        ssoUrl = el.getAttribute("Location") ?? undefined;
        break;
      }
    }

    const certs = xml.getElementsByTagNameNS(dsNS, "X509Certificate");
    const certificateString =
      certs.length > 0 ? certs[0].textContent?.trim() : undefined;
    const certificate = `-----BEGIN CERTIFICATE-----\n${certificateString}\n-----END CERTIFICATE-----`;

    setOktaMetadata({
      idpEntityId: entityID,
      idpRedirectUrl: ssoUrl,
      idpX509Certificate: certificate,
    } as OktaMetadata);
  }

  return (
    <>
      <div className="space-y-4 text-sm">
        <p className="font-medium">Create your Okta SAML application:</p>
        <ol className="list-decimal list-inside space-y-2">
          <li>Click on the "Sign On" tab</li>
          <li>Scroll down to where you see "Metadata URL"</li>
          <li>Click "Copy"</li>
          <li>Paste the SAML metadata URL below.</li>
        </ol>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)}>
          <div className="space-y-6">
            <FormField
              control={form.control}
              name="oktaMetadataUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Okta Metadata URL</FormLabel>
                  <FormDescription></FormDescription>
                  <FormMessage />
                  <FormControl>
                    <Input
                      type="url"
                      placeholder=""
                      {...field}
                      onChange={handleUrlChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                disabled={!form.formState.isDirty}
                type="submit"
                size="sm"
              >
                Continue
              </Button>
            </DialogFooter>
          </div>
        </form>
      </Form>
    </>
  );
}
