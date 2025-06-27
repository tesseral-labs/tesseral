import { useMutation, useQuery } from "@connectrpc/connect-query";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { consoleSearch } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleSearchResponse } from "@/gen/tesseral/backend/v1/backend_pb";
import {
  APIKey,
  Organization,
  User,
} from "@/gen/tesseral/backend/v1/models_pb";
import { useGlobalSearch } from "@/lib/search";

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

  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [query, setQuery] = useState<string>("");
  const [users, setUsers] = useState<User[]>([]);

  const searchMutation = useMutation(consoleSearch);

  useEffect(() => {
    if (!query || query.length < 3) {
      setApiKeys([]);
      setOrganizations([]);
      setUsers([]);
      return;
    }

    async function search() {
      try {
        const { apiKeys, organizations, users } =
          await searchMutation.mutateAsync({
            query,
          });
        setApiKeys(apiKeys);
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
            <CommandEmpty>No results found.</CommandEmpty>
            {(apiKeys?.length || 0) > 0 && (
              <CommandGroup heading="API Keys">
                {apiKeys?.map((apiKey) => (
                  <CommandItem
                    key={apiKey.id}
                    onSelect={() => {
                      setOpen(false);
                      setQuery("");
                      navigate(
                        `/organizations/${apiKey.organizationId}/api-keys/${apiKey.id}`,
                      );
                    }}
                  >
                    {apiKey.displayName} ({apiKey.id})
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
            {(organizations?.length || 0) > 0 && (
              <CommandGroup heading="Organizations">
                {organizations?.map((org) => (
                  <CommandItem
                    key={org.id}
                    onSelect={() => {
                      setOpen(false);
                      setQuery("");
                      navigate(`/organizations/${org.id}`);
                    }}
                  >
                    {org.displayName} ({org.id})
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
            {(users?.length || 0) > 0 && (
              <CommandGroup heading="Users">
                {users?.map((user) => (
                  <CommandItem
                    key={user.id}
                    onSelect={() => {
                      setOpen(false);
                      setQuery("");
                      navigate(
                        `/organizations/${user.organizationId}/users/${user.id}`,
                      );
                    }}
                  >
                    {user.displayName && <span>{user.displayName}</span>}{" "}
                    {user.email} ({user.id})
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
