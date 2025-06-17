import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { format } from "date-fns";
import {
  CalendarIcon,
  ChevronDown,
  ChevronRight,
  ExternalLink,
  Filter,
  Tag,
  XIcon,
} from "lucide-react";
import React, { useEffect, useMemo, useState } from "react";
import { DateRange } from "react-day-picker";
import { Link } from "react-router";

import {
  consoleListAuditLogEventNames,
  consoleListAuditLogEvents,
  getAPIKey,
  getBackendAPIKey,
  getSession,
  getUser,
} from "@/gen/tesseral/backend/v1/backend-BackendService_connectquery";
import { ConsoleListAuditLogEventsRequest } from "@/gen/tesseral/backend/v1/backend_pb";
import {
  APIKey,
  AuditLogEvent,
  BackendAPIKey,
  Session,
  User,
} from "@/gen/tesseral/backend/v1/models_pb";
import { cn } from "@/lib/utils";

import { ValueCopier } from "../core/ValueCopier";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { Calendar } from "../ui/calendar";
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "../ui/hover-card";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Separator } from "../ui/separator";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

export function ListAuditLogEventsTable({
  listParams,
}: {
  listParams: ConsoleListAuditLogEventsRequest;
}) {
  const [date, setDate] = React.useState<DateRange | undefined>(undefined);
  const [eventName, setEventName] = useState<string>("");
  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});

  const stableListParams = useMemo(() => {
    const params = {
      ...listParams,
      pageToken: "",
    } as ConsoleListAuditLogEventsRequest;

    // Set the event name filter if provided
    if (eventName && eventName.length > 0) {
      params.filterEventName = eventName;
    }

    // Set date filters if date range is selected
    if (date?.from) {
      params.filterStartTime = timestampFromDate(date.from);
    }
    if (date?.to) {
      const end = new Date(date.to);
      end.setHours(23, 59, 59, 999);
      params.filterEndTime = timestampFromDate(end);
    }

    return params;
  }, [date, eventName, listParams]);

  const {
    data: listAuditLogEventsResponse,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery(consoleListAuditLogEvents, stableListParams, {
    pageParamKey: "pageToken",
    getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
  });
  const { data: listAuditLogEventNamesResponse } = useQuery(
    consoleListAuditLogEventNames,
    {
      actorApiKeyId: listParams.actorApiKeyId,
      actorBackendApiKeyId: listParams.actorBackendApiKeyId,
      actorSessionId: listParams.actorSessionId,
      actorUserId: listParams.actorUserId,
      organizationId: listParams.organizationId,
      resourceType: listParams.resourceType,
    },
  );

  const auditLogEvents =
    listAuditLogEventsResponse?.pages.flatMap((page) => page.auditLogEvents) ||
    [];

  function toggleRow(eventId: string) {
    setExpandedRows((prev) => ({
      ...prev,
      [eventId]: !prev[eventId],
    }));
  }

  return (
    <>
      <div className="px-4 py-2 border-y bg-muted/40 space-y-2">
        <div className="font-semibold text-sm flex items-center gap-2">
          <Filter className="w-3 h-3" />
          <span>Filters</span>
        </div>
        <div className="flex items-center gap-4">
          <div className="space-y-1">
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant={"outline"}
                  className="justify-start text-left font-normal"
                >
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {date?.from ? (
                    date.to ? (
                      <>
                        {format(date.from, "LLL dd, y")} -{" "}
                        {format(date.to, "LLL dd, y")}
                      </>
                    ) : (
                      format(date.from, "LLL dd, y")
                    )
                  ) : (
                    <span className="text-muted-foreground">
                      Pick a date range
                    </span>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="range"
                  selected={date}
                  onSelect={setDate}
                  numberOfMonths={1}
                />
                {date && (date.to || date.from) && (
                  <div className="p-2">
                    <Button
                      variant="outline"
                      size="sm"
                      className="w-full"
                      onClick={() => setDate(undefined)}
                    >
                      Clear Date Range
                    </Button>
                  </div>
                )}
              </PopoverContent>
            </Popover>
          </div>

          <div className="space-y-2">
            <Select value={eventName} onValueChange={setEventName}>
              <SelectTrigger className="bg-white hover:bg-muted">
                <Tag />
                <SelectValue placeholder="Pick an event" />
              </SelectTrigger>
              <SelectContent>
                {eventName && eventName.length > 0 && (
                  <Button
                    className="w-full"
                    variant="outline"
                    size="sm"
                    onClick={() => setEventName("")}
                  >
                    Clear Filter
                  </Button>
                )}
                <SelectGroup>
                  <SelectLabel>Event</SelectLabel>
                  {listAuditLogEventNamesResponse?.eventNames.map(
                    (eventName) => (
                      <SelectItem key={eventName} value={eventName}>
                        {eventName}
                      </SelectItem>
                    ),
                  )}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          {((date && date.from) || (eventName && eventName.length > 0)) && (
            <div className="ml-auto">
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  setDate(undefined);
                  setEventName("");
                }}
              >
                <XIcon />
                Clear Filters
              </Button>
            </div>
          )}
        </div>
      </div>
      {auditLogEvents.length === 0 ? (
        <div className="py-4 text-sm text-center text-muted-foreground mt-8">
          No log events found.
        </div>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead></TableHead>
              <TableHead>Event</TableHead>
              <TableHead>Actor</TableHead>
              <TableHead className="text-right">Time</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {auditLogEvents.map((event) => (
              <>
                <AuditLogEventRow
                  key={event.id}
                  event={event}
                  expandedRows={expandedRows}
                  toggleRow={toggleRow}
                />
                {expandedRows[event.id] && (
                  <AuditLogEventDetails event={event} />
                )}
              </>
            ))}
          </TableBody>
        </Table>
      )}

      {hasNextPage && (
        <div className="flex justify-center mt-8">
          <Button
            variant="outline"
            size="sm"
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            Load More
          </Button>
        </div>
      )}
    </>
  );
}

function AuditLogEventRow({
  event,
  expandedRows,
  toggleRow,
}: {
  event: AuditLogEvent;
  expandedRows: Record<string, boolean>;
  toggleRow: (eventId: string) => void;
}) {
  function parseActor(event: AuditLogEvent): string {
    if (event.actorUserId) {
      return event.actorUserId;
    }
    if (event.actorApiKeyId) {
      return event.actorApiKeyId;
    }
    if (event.actorBackendApiKeyId) {
      return event.actorBackendApiKeyId;
    }

    return "System";
  }

  return (
    <TableRow onClick={() => toggleRow(event.id)} className="cursor-pointer">
      <TableCell className="align-middle">
        {expandedRows[event.id] ? (
          <ChevronDown className="h-4 w-4" />
        ) : (
          <ChevronRight className="h-4 w-4" />
        )}
      </TableCell>
      <TableCell>
        <span className="font-mono text-xs bg-muted px-2 py-1">
          {event.eventName}
        </span>
      </TableCell>
      <TableCell>
        <AuditLogEventActor event={event} />
      </TableCell>
      <TableCell className="text-right">
        <span className="text-muted-foreground">
          {event.eventTime && format(timestampDate(event.eventTime), "PPpp")}
        </span>
      </TableCell>
    </TableRow>
  );
}

function AuditLogEventActor({
  event: auditLogEvent,
}: {
  event: AuditLogEvent;
}) {
  const getApiKeyMutation = useMutation(getAPIKey);
  const getBackendAPIKeyMutation = useMutation(getBackendAPIKey);
  const getSessionMutation = useMutation(getSession);
  const getUserMutation = useMutation(getUser);

  const [apiKeyActor, setApiKeyActor] = useState<Record<string, any>>();
  const [sessionActor, setSessionActor] = useState<Record<string, any>>();
  const [userActor, setUserActor] = useState<Record<string, any>>();

  useEffect(() => {
    if (auditLogEvent) {
      (async () => {
        if (auditLogEvent.actorApiKeyId) {
          const { apiKey } = await getApiKeyMutation.mutateAsync({
            id: auditLogEvent.actorApiKeyId,
          });
          if (!apiKey) return;

          const apiKeyActor = {
            ...apiKey,
            createTime: timestampDate(apiKey.createTime!).toISOString(),
            updateTime: timestampDate(apiKey.updateTime!).toISOString(),
          };

          delete (apiKeyActor as any).$typeName;
          setApiKeyActor(apiKeyActor);
        }
        if (auditLogEvent.actorBackendApiKeyId) {
          const { backendApiKey } = await getBackendAPIKeyMutation.mutateAsync({
            id: auditLogEvent.actorBackendApiKeyId,
          });
          if (!backendApiKey) return;

          const backendApiKeyActor = {
            ...backendApiKey,
            createTime: timestampDate(backendApiKey.createTime!).toISOString(),
            updateTime: timestampDate(backendApiKey.updateTime!).toISOString(),
          };
          delete (backendApiKeyActor as any).$typeName;
          setApiKeyActor(backendApiKeyActor);
        }
        if (auditLogEvent.actorSessionId) {
          const { session } = await getSessionMutation.mutateAsync({
            id: auditLogEvent.actorSessionId,
          });
          if (!session) return;
          const sessionActor = {
            ...session,
            createTime: timestampDate(session.createTime!).toISOString(),
            expireTime: timestampDate(session.expireTime!).toISOString(),
            lastActiveTime: timestampDate(
              session.lastActiveTime!,
            ).toISOString(),
          };
          delete (sessionActor as any).$typeName;
          setSessionActor(session);

          const { user } = await getUserMutation.mutateAsync({
            id: session.userId,
          });
          if (!user) return;
          const userActor = {
            ...user,
            createTime: timestampDate(user.createTime!).toISOString(),
            updateTime: timestampDate(user.updateTime!).toISOString(),
          };
          delete (userActor as any).$typeName;
          setUserActor(userActor);
        }
        if (auditLogEvent.actorUserId) {
          const { user } = await getUserMutation.mutateAsync({
            id: auditLogEvent.actorUserId,
          });
          if (!user) return;
          const userActor = {
            ...user,
            createTime: timestampDate(user.createTime!).toISOString(),
            updateTime: timestampDate(user.updateTime!).toISOString(),
          };
          delete (userActor as any).$typeName;
          setUserActor(userActor);
        }
      })();
    }
  }, [auditLogEvent]);

  return (
    <>
      {apiKeyActor || sessionActor || userActor ? (
        <HoverCard>
          <HoverCardTrigger>
            {apiKeyActor ? (
              <Badge variant="secondary">
                <span className="font-mono">
                  {apiKeyActor.displayName || apiKeyActor.id}
                </span>
              </Badge>
            ) : sessionActor ? (
              <Badge variant="secondary">
                <span className="font-mono">
                  {userActor && userActor.email}
                </span>
              </Badge>
            ) : userActor ? (
              <Badge variant="secondary">
                <span className="font-mono">{userActor.email}</span>
              </Badge>
            ) : (
              <></>
            )}
          </HoverCardTrigger>
          <HoverCardContent>
            {apiKeyActor && (
              <>
                <div className="font-semibold">API Key</div>
                <pre className="text-xs">
                  {JSON.stringify(apiKeyActor, null, 2)}
                </pre>
              </>
            )}
            {userActor && (
              <>
                <div className="font-semibold">User</div>
                <pre className="text-xs">
                  {JSON.stringify(userActor, null, 2)}
                </pre>
              </>
            )}
          </HoverCardContent>
        </HoverCard>
      ) : (
        <Badge variant="outline">
          <span className="font-mono">System</span>
        </Badge>
      )}
    </>
  );
}

export function AuditLogEventDetails({ event }: { event: AuditLogEvent }) {
  return (
    <TableRow className="bg-muted/40">
      <TableCell colSpan={4} className="p-4">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="w-full md:w-1/2">
            <AuditLogEventActorDetails event={event} />
          </div>

          <Separator orientation="vertical" />

          <div className="w-full md:w-1/2">
            <div className="space-y-4">
              <div className="space-y-2">
                <div className="text-base font-semibold">Event</div>
                <div className="text-xs font-mono">
                  <ValueCopier value={event.id} label="Audit Log Event ID" />
                </div>
              </div>
              <div className="space-y-2">
                <div className="text-base font-semibold">Details</div>
                <div className="font-mono text-xs whitespace-pre-wrap break-all p-2 rounded bg-white border">
                  {JSON.stringify(event.eventDetails, null, 2)}
                </div>
              </div>
            </div>
          </div>
        </div>
      </TableCell>
    </TableRow>
  );
}

function AuditLogEventActorDetails({ event }: { event: AuditLogEvent }) {
  const { actorApiKeyId, actorBackendApiKeyId, actorSessionId, actorUserId } =
    event;

  const getApiKeyMutation = useMutation(getAPIKey);
  const getBackendApiKeyMutation = useMutation(getBackendAPIKey);
  const getSessionMutation = useMutation(getSession);
  const getUserMutation = useMutation(getUser);

  const [apiKey, setApiKey] = useState<APIKey>();
  const [backendApiKey, setBackendApiKey] = useState<BackendAPIKey>();
  const [session, setSession] = useState<Session>();
  const [user, setUser] = useState<User>();

  useEffect(() => {
    if (actorApiKeyId && !apiKey && !getApiKeyMutation.isPending) {
      getApiKeyMutation
        .mutateAsync({ id: actorApiKeyId })
        .then((response) => response.apiKey)
        .then((apiKey) => setApiKey(apiKey));
    }
  }, [actorApiKeyId, apiKey, getApiKeyMutation]);

  useEffect(() => {
    if (
      actorBackendApiKeyId &&
      !backendApiKey &&
      !getBackendApiKeyMutation.isPending
    ) {
      getBackendApiKeyMutation
        .mutateAsync({ id: actorBackendApiKeyId })
        .then((response) => response.backendApiKey)
        .then((backendApiKey) => setBackendApiKey(backendApiKey));
    }
  }, [actorBackendApiKeyId, backendApiKey, getBackendApiKeyMutation]);

  useEffect(() => {
    if (actorSessionId && !session && !getSessionMutation.isPending) {
      getSessionMutation
        .mutateAsync({ id: actorSessionId })
        .then((response) => response.session)
        .then((session) => setSession(session));
    }
  }, [actorSessionId, getSessionMutation, session]);

  useEffect(() => {
    if (actorUserId && !user && !getUserMutation.isPending) {
      getUserMutation
        .mutateAsync({ id: actorUserId })
        .then((response) => response.user)
        .then((user) => setUser(user));
    }
  }, [actorUserId, getUserMutation, user]);

  useEffect(() => {
    if (session && session.userId && !user && !getUserMutation.isPending) {
      getUserMutation
        .mutateAsync({ id: session.userId })
        .then((response) => response.user)
        .then((user) => setUser(user));
    }
  }, [getUserMutation, session, user]);

  return (
    <div className="space-y-4">
      {apiKey && (
        <div className="space-y-1">
          <div className="font-semibold text-base">API Key</div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/organizations/${apiKey.organizationId}/api-keys/${apiKey.id}`}
          >
            {apiKey.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}
      {backendApiKey && (
        <div className="space-y-1">
          <div className="font-semibold text-base">Backend API Key</div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/settings/api-keys/${backendApiKey.id}`}
          >
            {backendApiKey.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}
      {user && (
        <div className="space-y-1">
          <div className="font-semibold text-base">User</div>
          {user.displayName && (
            <div className="font-medium">{user.displayName}</div>
          )}
          <div
            className={cn(
              user.displayName ? "text-muted-foreground" : "font-medium",
            )}
          >
            {user.email}
          </div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/organizations/${user.organizationId}/users/${user.id}`}
          >
            {user.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}
      {session && user && (
        <div className="space-y-1">
          <div className="font-semibold text-base">Session</div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/organizations/${user.organizationId}/users/${session.userId}/sessions/${session.id}`}
          >
            {session.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}
      {!apiKey && !session && !user && (
        <div className="text-muted-foreground text-sm">System</div>
      )}
    </div>
  );
}
