import { zodResolver } from "@hookform/resolvers/zod";
import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";

const schema = z.object({});

export function CreateSAMLConnectionButton() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {},
  });

  async function handleSubmit(data: z.infer<typeof schema>) {
    // Handle form submission logic here
    console.log("Form submitted with data:", data);
    // You can call an API to create the SAML connection
  }

  return (
    <AlertDialog>
      <AlertDialogTrigger>
        <Button variant="outline">Create SAML Connection</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create SAML Connection</AlertDialogTitle>
          <AlertDialogDescription>
            To create a SAML connection, you will need to provide the necessary
            details such as the SAML metadata URL or file, and any additional
            configuration required by your Identity Provider (IdP). Please
            ensure that you have the required information ready before
            proceeding.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          ></form>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <Button type="submit">Create Connection</Button>
          </AlertDialogFooter>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
}
