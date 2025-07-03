import { useMutation, useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { ChangeEvent, useRef, useState } from "react";
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
import { Separator } from "@/components/ui/separator";
import {
  getSAMLConnection,
  updateSAMLConnection,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { SAMLMetadata, parseSamlMetadata } from "@/lib/saml";

const schema = z.object({
  googleSamlMetadata: z.string(),
});

export function DownloadGoogleSamlMetadata() {
  const filePickerRef = useRef<HTMLInputElement>(null);
  const { samlConnectionId } = useParams();
  const navigate = useNavigate();

  const { refetch } = useQuery(getSAMLConnection, {
    id: samlConnectionId,
  });
  const updateSamlConnectionMutation = useMutation(updateSAMLConnection);

  const [metadata, setMetadata] = useState<SAMLMetadata>();

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      googleSamlMetadata: "",
    },
  });

  function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
    if (!e.target.files || e.target.files.length === 0) {
      form.setValue("googleSamlMetadata", "");
      return;
    }

    const file = e.target.files[0];
    if (file.type !== "application/xml" && file.type !== "text/xml") {
      form.setError("googleSamlMetadata", {
        type: "manual",
        message: "Please upload a valid XML file.",
      });
      return;
    }

    const reader = new FileReader();
    form.setValue("googleSamlMetadata", file.name, {
      shouldDirty: true,
      shouldTouch: true,
      shouldValidate: true,
    });

    reader.onload = (event) => {
      const xmlString = event.target?.result as string;
      try {
        try {
          const metadata = parseSamlMetadata(xmlString);
          setMetadata(metadata);
        } catch {
          form.setError("googleSamlMetadata", {
            type: "manual",
            message: "Error parsing XML file.",
          });
        }
      } catch {
        form.setError("googleSamlMetadata", {
          type: "manual",
          message: "Error parsing XML file.",
        });
      }
    };
    reader.onerror = () => {
      form.setError("googleSamlMetadata", {
        type: "manual",
        message: "Error reading file.",
      });
    };
    reader.readAsText(file);
  }

  async function handleSubmit() {
    try {
      await updateSamlConnectionMutation.mutateAsync({
        id: samlConnectionId,
        samlConnection: {
          idpEntityId: metadata?.idpEntityId,
          idpRedirectUrl: metadata?.idpRedirectUrl,
          idpX509Certificate: metadata?.idpX509Certificate,
        },
      });
      await refetch();
      form.reset();
      toast.success("SAML connection updated successfully.");
      navigate(
        `/organization/saml-connections/${samlConnectionId}/setup/google/configure`,
      );
    } catch {
      toast.error("Failed to update SAML connection. Please try again.");
      return;
    }
  }

  return (
    <>
      <div className="space-y-4 text-sm">
        <img
          className="rounded-xl max-w-full border shadow-md"
          src="/videos/saml-setup-wizard/google/metadata.gif"
        />

        <p className="font-medium">Download IdP metadata:</p>

        <ol className="list-decimal list-inside space-y-2">
          <li>Click on "Download Metadata".</li>
          <li>
            Locate the downloaded <code>GoogleIDPMetadata.xml</code> file.
          </li>
          <li>
            Upload the <code>GoogleIDPMetadata.xml</code> file here.
          </li>
        </ol>
      </div>

      <Separator />

      <Form {...form}>
        <form className="mt-4" onSubmit={form.handleSubmit(handleSubmit)}>
          <FormField
            control={form.control}
            name="googleSamlMetadata"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Google IDP Metadata</FormLabel>
                <FormDescription></FormDescription>
                <FormMessage />
                <FormControl>
                  <div className="flex items-center space-x-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => filePickerRef.current?.click()}
                    >
                      Choose file
                    </Button>
                    <div className="text-sm">
                      {field.value || "No file chosen"}
                    </div>

                    <input
                      accept=".xml"
                      className="hidden"
                      onChange={handleFileChange}
                      ref={filePickerRef}
                      type="file"
                    />
                  </div>
                </FormControl>
              </FormItem>
            )}
          />

          <DialogFooter>
            <Button disabled={!form.formState.isDirty} size="sm">
              Continue
            </Button>
          </DialogFooter>
        </form>
      </Form>
    </>
  );
}
