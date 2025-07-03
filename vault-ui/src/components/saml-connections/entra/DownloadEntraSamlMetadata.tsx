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
  entraSamlMetadata: z.string(),
});

export function DownloadEntraSamlMetadata() {
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
      entraSamlMetadata: "",
    },
  });

  function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
    if (!e.target.files || e.target.files.length === 0) {
      form.setValue("entraSamlMetadata", "");
      return;
    }

    const file = e.target.files[0];
    if (file.type !== "application/xml" && file.type !== "text/xml") {
      form.setError("entraSamlMetadata", {
        type: "manual",
        message: "Please upload a valid XML file.",
      });
      return;
    }

    const reader = new FileReader();
    form.setValue("entraSamlMetadata", file.name, {
      shouldDirty: true,
      shouldTouch: true,
      shouldValidate: true,
    });

    reader.onload = (event) => {
      const xmlString = event.target?.result as string;
      try {
        const metadata = parseSamlMetadata(xmlString);
        setMetadata(metadata);
      } catch {
        form.setError("entraSamlMetadata", {
          type: "manual",
          message: "Error parsing XML file.",
        });
      }
    };
    reader.onerror = () => {
      form.setError("entraSamlMetadata", {
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
        `/organization/saml-connections/${samlConnectionId}/setup/entra/users`,
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
          src="/videos/saml-setup-wizard/entra/metadata.gif"
        />

        <p className="font-medium">
          Download your application's Federation Metadata XML:
        </p>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Close the "Basic SAML Connection" dialog if you haven't already
          </li>
          <li>
            Back in the "Single sign-On" section scroll down to the "SAML
            Certificates" section
          </li>
          <li>
            Click on the "Download" link next to "Federation Metadata XML"
          </li>
          <li>Your browser now downloads a file</li>
          <li>Upload that file below</li>
        </ol>

        <Separator />

        <Form {...form}>
          <form className="mt-4" onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="entraSamlMetadata"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Federated Metadata</FormLabel>
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
      </div>
    </>
  );
}
