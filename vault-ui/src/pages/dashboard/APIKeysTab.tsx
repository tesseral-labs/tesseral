import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { format } from "date-fns";
import { CalendarIcon, CirclePlus, Copy, LoaderCircle } from "lucide-react";
import { DateTime } from "luxon";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { z } from "zod";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  createAPIKey,
  getOrganization,
  getProject,
  listAPIKeys,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { APIKey } from "@/gen/tesseral/frontend/v1/models_pb";
import { cn } from "@/lib/utils";

export function APIKeysTab() {
  const { data: getOrganizationResponse } = useQuery(getOrganization);
  const { data: getProjectResponse } = useQuery(getProject);
  const {
    data: listApiKeysResponses,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(
    listAPIKeys,
    {
      organizationId: getOrganizationResponse?.organization?.id,
      pageToken: "",
    },
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const apiKeys = listApiKeysResponses?.pages.flatMap((page) => page.apiKeys);

  return (
    <Card>
      <CardHeader className="flex flex-row justify-between gap-x-4">
        <div>
          <CardTitle>API Keys</CardTitle>
          <CardDescription>
            API keys are used to authenticate and authorize access to the API.
            You can create, manage, and revoke API keys for your organization.
          </CardDescription>
        </div>

        <CreateAPIKeyButton />
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Display Name</TableHead>
              <TableHead>ID</TableHead>
              <TableHead>Value</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Expires</TableHead>
              <TableHead>Created At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {apiKeys &&
              apiKeys.map((apiKey) => (
                <TableRow key={apiKey.id}>
                  <TableCell>
                    <Link to={`/organization-settings/api-keys/${apiKey.id}`}>
                      {apiKey.displayName}
                    </Link>
                  </TableCell>
                  <TableCell>
                    <Link to={`/organization-settings/api-keys/${apiKey.id}`}>
                      {apiKey.id}
                    </Link>
                  </TableCell>
                  <TableCell>
                    {apiKey.secretTokenSuffix ? (
                      <span className="font-mono text-sm">
                        {getProjectResponse?.project?.apiKeySecretTokenPrefix ||
                          "api_key_"}
                        ...{apiKey.secretTokenSuffix}
                      </span>
                    ) : (
                      "—"
                    )}
                  </TableCell>
                  <TableCell>
                    {apiKey.revoked ? (
                      <span>Active</span>
                    ) : (
                      <span>Revoked</span>
                    )}
                  </TableCell>
                  <TableCell>
                    {apiKey.expireTime
                      ? DateTime.fromJSDate(
                          timestampDate(apiKey.expireTime),
                        ).toRelative()
                      : "never"}
                  </TableCell>
                  <TableCell>
                    {apiKey.createTime &&
                      DateTime.fromJSDate(
                        timestampDate(apiKey.createTime),
                      ).toRelative()}
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>

        {hasNextPage && (
          <div className="flex justify-center mt-8">
            <Button
              className="mt-4"
              variant="outline"
              onClick={() => fetchNextPage()}
            >
              {isFetchingNextPage && (
                <LoaderCircle className="h-4 w-4 animate-spin" />
              )}
              Load more
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

const schema = z.object({
  displayName: z.string(),
  expireTime: z.string(),
});

function CreateAPIKeyButton() {
  const [apiKey, setApiKey] = useState<APIKey>();

  const [customDate, setCustomDate] = useState<Date>();

  const { organizationId } = useParams();
  const { refetch } = useQuery(listAPIKeys, {
    organizationId,
  });

  const createApiKeyMutation = useMutation(createAPIKey);

  const form = useForm<z.infer<typeof schema>>({
    defaultValues: {
      displayName: "",
      expireTime: "1 day",
    },
  });

  const handleSubmit = async (data: z.infer<typeof schema>) => {
    const createParams: Record<string, any> = {
      organizationId: organizationId!,
      displayName: data.displayName,
    };

    switch (data.expireTime) {
      case "1 day":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 24 * 60 * 60 * 1000),
        );
        break;
      case "7 days":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
        );
        break;
      case "30 days":
        createParams.expireTime = timestampFromDate(
          new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
        );
        break;
      case "custom":
        if (customDate) {
          createParams.expireTime = timestampFromDate(customDate);
        }
        break;
      case "noexpire":
        break;
    }

    const { apiKey } = await createApiKeyMutation.mutateAsync(createParams);

    if (apiKey) {
      setApiKey(apiKey);

      toast.success("API Key created successfully");

      await refetch();
    }
  };

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="outline">
          <CirclePlus className="h-4 w-4" />
          Create API Key
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create API Key</AlertDialogTitle>
        </AlertDialogHeader>

        {!apiKey ? (
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="displayName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Display Name</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="expireTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Expire time</FormLabel>
                    <FormControl>
                      <div className="flex flex-row gap-2">
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            field.onChange(value);

                            console.log(value);
                          }}
                        >
                          <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Pick a custom date" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="1 day">1 day</SelectItem>
                            <SelectItem value="7 days">7 days</SelectItem>
                            <SelectItem value="30 days">30 days</SelectItem>
                            <SelectItem value="custom">Custom</SelectItem>
                            <SelectItem value="noexpire">
                              No expiration
                            </SelectItem>
                          </SelectContent>
                        </Select>

                        {field.value === "custom" && (
                          <Popover>
                            <PopoverTrigger asChild>
                              <Button
                                variant={"outline"}
                                className={cn(
                                  "w-[270px] justify-start text-left font-normal",
                                  !customDate && "text-muted-foreground",
                                )}
                              >
                                <CalendarIcon className="mr-2 h-4 w-4" />
                                {customDate ? (
                                  format(customDate, "PPP")
                                ) : (
                                  <span>Pick a date</span>
                                )}
                              </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0">
                              <Calendar
                                mode="single"
                                selected={customDate}
                                onSelect={setCustomDate}
                                initialFocus
                              />
                            </PopoverContent>
                          </Popover>
                        )}
                      </div>
                    </FormControl>
                  </FormItem>
                )}
              />

              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <Button type="submit">Save</Button>
              </AlertDialogFooter>
            </form>
          </Form>
        ) : (
          <div className="space-y-4">
            <div className="text-muted-foreground text-sm text-wrap">
              This is your secret token. This token will not be shown again, so
              please save it now.
            </div>

            <div className="p-2 bg-muted text-muted-foreground font-mono text-xs overflow-x-hidden text-wrap word-break rounded">
              {apiKey.secretToken}
            </div>

            <Button
              variant="outline"
              onClick={() => {
                navigator.clipboard.writeText(apiKey.secretToken);
                toast.success("API Key copied to clipboard");
              }}
            >
              <Copy className="h-4 w-4" />
              Copy
            </Button>

            <AlertDialogFooter className="flex justify-end">
              <AlertDialogCancel>Done</AlertDialogCancel>
              <Link
                to={`/organizations/${organizationId}/api-keys/${apiKey.id}`}
              >
                <Button>Manage API Key</Button>
              </Link>
            </AlertDialogFooter>
          </div>
        )}
      </AlertDialogContent>
    </AlertDialog>
  );
}
