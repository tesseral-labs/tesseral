import { useMutation } from '@connectrpc/connect-query';
import { zodResolver } from '@hookform/resolvers/zod';
import { LoaderCircleIcon } from 'lucide-react';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { LoginFlowCard } from '@/components/login/LoginFlowCard';
import { Button } from '@/components/ui/button';
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  onboardingCreateProjects,
} from '@/gen/tesseral/intermediate/v1/intermediate-IntermediateService_connectquery';
import { useNavigate } from 'react-router';

const schema = z.object({
  displayName: z.string().nonempty(),
  productionAppUrl: z.string().url(),
  localhostAppUrl: z.string().url().startsWith('http://localhost:'),
});

export function CreateOrganizationPage() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      displayName: "",
      productionAppUrl: "",
      localhostAppUrl: "",
    },
  });

  const [submitting, setSubmitting] = useState(false);
  const { mutateAsync: onboardingCreateProjectsAsync } =
    useMutation(onboardingCreateProjects);
  const navigate = useNavigate();

  async function handleSubmit(values: z.infer<typeof schema>) {
    setSubmitting(true);

    await onboardingCreateProjectsAsync({
      displayName: values.displayName,
      prodUrl: values.productionAppUrl,
      devUrl: values.localhostAppUrl,
    });

    navigate('/');
  }

  return (
    <LoginFlowCard>
      <CardHeader>
        <CardTitle>Create a new project</CardTitle>
        <CardDescription>
          To get going quickly, we'll ask a few questions. All your answers here
          can be changed later.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form className="space-y-4" onSubmit={form.handleSubmit(handleSubmit)}>
            <FormField
              control={form.control}
              name="displayName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Company Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Example Corporation" {...field} />
                  </FormControl>
                  <FormDescription>
                    What's your company's name?
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="productionAppUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Production App URL</FormLabel>
                  <FormControl>
                    <Input placeholder="https://app.example.com" {...field} />
                  </FormControl>
                  <FormDescription>
                    Where does your app live in production?
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="localhostAppUrl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Development App URL</FormLabel>
                  <FormControl>
                    <Input placeholder="http://localhost:3000" {...field} />
                  </FormControl>
                  <FormDescription>
                    What localhost URL does your app live on?
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="mt-4 w-full" disabled={submitting}>
              {submitting && (
                <LoaderCircleIcon className="h-4 w-4 animate-spin" />
              )}
              Create Project
            </Button>
          </form>
        </Form>
      </CardContent>
    </LoginFlowCard>
  );
}
