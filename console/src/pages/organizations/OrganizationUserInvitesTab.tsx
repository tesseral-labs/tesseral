import React, { useState } from 'react';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import {
  createUserInvite,
  listUserInvites,
} from '@/gen/tesseral/backend/v1/backend-BackendService_connectquery';
import {
  ConsoleCard,
  ConsoleCardDetails,
  ConsoleCardDescription,
  ConsoleCardHeader,
  ConsoleCardTitle,
  ConsoleCardTableContent,
} from '@/components/ui/console-card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Link } from 'react-router-dom';
import { DateTime } from 'luxon';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { useNavigate, useParams } from 'react-router';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
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
import { Switch } from '@/components/ui/switch';
import { toast } from 'sonner';

export const OrganizationUserInvitesTab = () => {
  const { organizationId } = useParams();
  const { data: listUserInvitesResponse } = useQuery(listUserInvites, {
    organizationId,
  });

  return (
    <ConsoleCard>
      <ConsoleCardHeader>
        <ConsoleCardDetails>
          <ConsoleCardTitle>User Invites</ConsoleCardTitle>
          <ConsoleCardDescription>
            A User Invite lets outside collaborators join an Organization.
          </ConsoleCardDescription>
        </ConsoleCardDetails>
        <CreateUserInviteButton />
      </ConsoleCardHeader>
      <ConsoleCardTableContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead>Updated At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {listUserInvitesResponse?.userInvites?.map((userInvite) => (
              <TableRow key={userInvite.id}>
                <TableCell>
                  <Link
                    className="font-medium underline underline-offset-2 decoration-muted-foreground/40"
                    to={`/organizations/${organizationId}/user-invites/${userInvite.id}`}
                  >
                    {userInvite.email}
                  </Link>

                  {userInvite.owner && (
                    <Badge variant="outline" className="ml-2">
                      Owner
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="font-mono">{userInvite.id}</TableCell>
                <TableCell>
                  {userInvite.createTime &&
                    DateTime.fromJSDate(
                      timestampDate(userInvite.createTime),
                    ).toRelative()}
                </TableCell>
                <TableCell>
                  {userInvite.updateTime &&
                    DateTime.fromJSDate(
                      timestampDate(userInvite.updateTime),
                    ).toRelative()}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </ConsoleCardTableContent>
    </ConsoleCard>
  );
};

const schema = z.object({
  email: z.string().email(),
  owner: z.boolean(),
  sendEmail: z.boolean(),
});

const CreateUserInviteButton = () => {
  const { organizationId } = useParams();
  const createUserInviteMutation = useMutation(createUserInvite);
  /* eslint-disable @typescript-eslint/no-unsafe-call */
  // Currently there's an issue with the types of react-hook-form and zod
  // preventing the compiler from inferring the correct types.
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
      owner: false,
      sendEmail: true,
    },
  });
  /* eslint-enable @typescript-eslint/no-unsafe-call */
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  const handleSubmit = async (values: z.infer<typeof schema>) => {
    const { userInvite } = await createUserInviteMutation.mutateAsync({
      userInvite: {
        organizationId,
        email: values.email,
        owner: values.owner,
      },
      sendEmail: values.sendEmail,
    });

    toast.success('User Invite created');
    navigate(`/organizations/${organizationId}/user-invites/${userInvite?.id}`);
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger>
        <Button variant="outline">Create</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create User Invite</AlertDialogTitle>
          <AlertDialogDescription>
            A User Invite lets outside collaborators join an organization.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          {/* eslint-disable @typescript-eslint/no-unsafe-call */}
          {/* Currently there's an issue with the types of react-hook-form and zod 
          preventing the compiler from inferring the correct types.*/}
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-8"
          >
            {/* eslint-enable @typescript-eslint/no-unsafe-call */}
            <FormField
              control={form.control}
              name="email"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input className="max-w-96" {...field} />
                  </FormControl>
                  <FormDescription>
                    The outside collaborator's email. The collaborator will need
                    to verify this email before being able to join the
                    organization.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="owner"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Invite as owner</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Whether the collaborator will join as an owner.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="sendEmail"
              render={({ field }: { field: any }) => (
                <FormItem>
                  <FormLabel>Send email</FormLabel>
                  <FormControl>
                    <Switch
                      className="block"
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <FormDescription>
                    Send the collaborator an email to accept the invite.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <AlertDialogFooter className="mt-8">
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <Button type="submit">Create User Invite</Button>
            </AlertDialogFooter>
          </form>
        </Form>
      </AlertDialogContent>
    </AlertDialog>
  );
};
