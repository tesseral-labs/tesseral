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
  samlMetadata: z.string().min(1, "Please upload a valid XML file."),
});

export function DownloadOtherSamlMetadata() {
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
      samlMetadata: "",
    },
  });

  function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
    if (!e.target.files || e.target.files.length === 0) {
      form.setValue("samlMetadata", "");
      return;
    }

    const file = e.target.files[0];
    if (file.type !== "application/xml" && file.type !== "text/xml") {
      form.setError("samlMetadata", {
        type: "manual",
        message: "Please upload a valid XML file.",
      });
      return;
    }

    const reader = new FileReader();
    form.setValue("samlMetadata", file.name, {
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
        form.setError("samlMetadata", {
          type: "manual",
          message: "Error parsing XML file.",
        });
      }
    };
    reader.onerror = () => {
      form.setError("samlMetadata", {
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
      toast.success("SAML Connection configured successfully.");
      navigate(
        `/organization/saml-connections/${samlConnectionId}/setup/other/users`,
      );
    } catch {
      toast.error("Failed to update SAML connection. Please try again later.");
    }
  }
  return (
    <div className="space-y-4 text-sm">
      <p className="font-medium">Download SAML Metadata</p>
      <p>
        Find your new SAML application's metadata XML file. This is a file your
        Identity Provider allows you to download on each of your SAML
        applications.
      </p>
      <p>Upload that file below.</p>

      <Separator />

      <Form {...form}>
        <form className="mt-4" onSubmit={form.handleSubmit(handleSubmit)}>
          <FormField
            control={form.control}
            name="samlMetadata"
            render={({ field }) => (
              <FormItem>
                <FormLabel>SAML Metadata</FormLabel>
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
  );
}
