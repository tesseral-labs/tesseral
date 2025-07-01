import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { Building2, Check, CornerDownLeft, UserIcon } from "lucide-react";
import { DateTime } from "luxon";
import React, { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import {
  consoleSearch,
  getOrganization,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleSearchResponse } from "@/gen/tesseral/backend/v1/backend_pb";
import {
  APIKey,
  Organization,
  User,
} from "@/gen/tesseral/backend/v1/models_pb";
import { useGlobalSearch } from "@/lib/search";

import { Badge } from "../ui/badge";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandList,
} from "../ui/command";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";

export function Search() {
  const { open, setOpen } = useGlobalSearch();
  const navigate = useNavigate();

  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [query, setQuery] = useState<string>("");
  const [users, setUsers] = useState<User[]>([]);

  const searchMutation = useMutation(consoleSearch);

  useEffect(() => {
    if (!query || query.length < 3) {
      setOrganizations([]);
      setUsers([]);
      return;
    }

    async function search() {
      try {
        const { organizations, users } = await searchMutation.mutateAsync({
          limit: 3,
          query,
        });
        setOrganizations(organizations);
        setUsers(users);
      } catch (error) {
        toast.error("Failed to perform search.");
        console.error(error);
      }
    }

    search();
  }, [query]);

  return (
    <Dialog open={open || false} onOpenChange={setOpen}>
      <DialogHeader>
        <DialogTitle>Search</DialogTitle>
      </DialogHeader>
      <DialogContent showCloseButton={false}>
        <Command>
          <Input
            className="focus:ring-0 focus-visible:ring-0"
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for a resource..."
            value={query}
          />

          <CommandList>
            {query.length >= 3 && (
              <CommandEmpty>No results found.</CommandEmpty>
            )}
            {(organizations?.length || 0) > 0 && (
              <CommandGroup heading="Organizations">
                {organizations?.map((org) => (
                  <CommandItem
                    key={org.id}
                    className="group"
                    onSelect={() => {
                      setOpen(false);
                      setQuery("");
                      navigate(`/organizations/${org.id}`);
                    }}
                  >
                    <OrganizationResult organization={org} />
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
            {(users?.length || 0) > 0 && (
              <CommandGroup heading="Users">
                {users?.map((user) => (
                  <CommandItem
                    key={user.id}
                    className="group"
                    onSelect={() => {
                      setOpen(false);
                      setQuery("");
                      navigate(
                        `/organizations/${user.organizationId}/users/${user.id}`,
                      );
                    }}
                  >
                    <UserResult user={user} />
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </CommandList>
        </Command>
      </DialogContent>
    </Dialog>
  );
}

function OrganizationResult({ organization }: { organization: Organization }) {
  const authFactors = useMemo(() => {
    const factors = [];
    if (organization.logInWithEmail) factors.push("Email");
    if (organization.logInWithPassword) factors.push("Password");
    if (organization.logInWithGoogle) factors.push("Google");
    if (organization.logInWithMicrosoft) factors.push("Microsoft");
    if (organization.logInWithGithub) factors.push("GitHub");
    if (organization.logInWithSaml) factors.push("SAML");

    if (factors.length > 1) {
      factors[factors.length - 1] = `and ${factors[factors.length - 1]}`;
    }

    return factors.length > 0
      ? `${factors.join(", ")} enabled.`
      : "No login methods enabled.";
  }, [organization]);
  return (
    <div className="flex items-center gap-2 w-full">
      <div className="rounded w-6 h-6 flex items-center justify-center shadow-sm bg-gradient-to-br from-gray-50 to-gray-200">
        <Building2 className="h-3 w-3 text-gray-400" />
      </div>
      <div className="space-y-1">
        <div>
          {organization.displayName ? (
            <div className="text-sm font-medium">
              {organization.displayName}
            </div>
          ) : (
            <div className="text-sm font-medium font-mono text-muted-foreground">
              {organization.id}
            </div>
          )}
        </div>
        <div className="text-xs text-muted-foreground">
          <span>
            Created{" "}
            {organization.createTime &&
              DateTime.fromJSDate(
                timestampDate(organization.createTime),
              ).toRelative()}
            .{" "}
          </span>
          <span>{authFactors}</span>
        </div>
        {/* <div className="space-y-1 space-x-1">
          {organization.logInWithEmail && (
            <Badge variant="outline">
              <Check /> Email
            </Badge>
          )}
          {organization.logInWithPassword && (
            <Badge variant="outline">
              <Check /> Password
            </Badge>
          )}
          {organization.logInWithGoogle && (
            <Badge variant="outline">
              <Check /> Google
            </Badge>
          )}
          {organization.logInWithMicrosoft && (
            <Badge variant="outline">
              <Check /> Microsoft
            </Badge>
          )}
          {organization.logInWithGithub && (
            <Badge variant="outline">
              <Check /> GitHub
            </Badge>
          )}
          {organization.logInWithSaml && (
            <Badge variant="outline">
              <Check /> SAML
            </Badge>
          )}
          {!organization.logInWithEmail &&
            !organization.logInWithPassword &&
            !organization.logInWithGoogle &&
            !organization.logInWithMicrosoft &&
            !organization.logInWithGithub &&
            !organization.logInWithSaml && (
              <Badge variant="outline">No login methods enabled</Badge>
            )}
        </div> */}
      </div>
      <div className="h-6 w-6 items-center justify-center ml-auto hidden group-data-[selected=true]:flex">
        <CornerDownLeft className="w-4 h-4 text-muted-foreground/30" />
      </div>
    </div>
  );
}

function UserResult({ user }: { user: User }) {
  const authFactors = useMemo(() => {
    const factors = [];
    if (user.email) factors.push("Email");
    if (user.googleUserId) factors.push("Google");
    if (user.microsoftUserId) factors.push("Microsoft");
    if (user.githubUserId) factors.push("GitHub");

    if (factors.length > 1) {
      factors[factors.length - 1] = `and ${factors[factors.length - 1]}`;
    }

    return factors.length > 0
      ? `${factors.join(", ")} registered.`
      : "No login methods registered.";
  }, [user]);
  const { data: getOrganizationResponse } = useQuery(getOrganization, {
    id: user.organizationId,
  });

  const organization = getOrganizationResponse?.organization;

  return (
    <div className="flex items-center gap-2 w-full">
      <div className="rounded w-6 h-6 flex items-center justify-center shadow-sm bg-gradient-to-br from-gray-50 to-gray-200">
        <UserIcon className="h-3 w-3 text-gray-400" />
      </div>
      <div className="space-y-1">
        {organization && (
          <>
            <div className="flex items-center text-sm font-medium gap-2">
              {user.email}
              <Badge className="text-xs" variant="outline">
                <Building2 />
                {organization.displayName}
              </Badge>
            </div>
            <div className="text-xs text-muted-foreground">
              <span>
                Created{" "}
                {user.createTime &&
                  DateTime.fromJSDate(
                    timestampDate(user.createTime),
                  ).toRelative()}
                .{" "}
              </span>
              <span>{authFactors}</span>
            </div>
          </>
        )}
      </div>
      <div className="h-6 w-6 items-center justify-center ml-auto hidden group-data-[selected=true]:flex">
        <CornerDownLeft className="w-4 h-4 text-muted-foreground/30" />
      </div>
    </div>
  );
}
