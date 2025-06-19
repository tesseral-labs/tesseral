import { useMutation, useQuery } from "@connectrpc/connect-query";
import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

import { consoleSearch } from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleSearchResponse } from "@/gen/tesseral/backend/v1/backend_pb";
import { useGlobalSearch } from "@/lib/search";

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "../ui/command";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";

export function Search() {
  const { open, setOpen } = useGlobalSearch();
  const navigate = useNavigate();

  const [query, setQuery] = useState<string>("");

  const { data: result } = useQuery(
    consoleSearch,
    {
      query,
    },
    {
      retry: false,
    },
  );

  return (
    <Dialog open={open || false} onOpenChange={setOpen}>
      <DialogHeader>
        <DialogTitle>Search</DialogTitle>
      </DialogHeader>
      <DialogContent>
        <Command>
          <Input
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for a resource..."
            value={query}
          />

          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            {(result?.apiKeys?.length || 0) > 0 && (
              <CommandGroup heading="API Keys">
                {result?.apiKeys?.map((apiKey) => (
                  <CommandItem
                    key={apiKey.id}
                    onSelect={() => {
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
            {(result?.organizations?.length || 0) > 0 && (
              <CommandGroup heading="Organizations">
                {result?.organizations?.map((org) => (
                  <CommandItem
                    key={org.id}
                    onSelect={() => {
                      navigate(`/organizations/${org.id}`);
                    }}
                  >
                    {org.displayName} ({org.id})
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
            {(result?.users?.length || 0) > 0 && (
              <CommandGroup heading="Users">
                {result?.users?.map((user) => (
                  <CommandItem
                    key={user.id}
                    onSelect={() => {
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
