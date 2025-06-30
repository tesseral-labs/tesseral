import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useQuery } from "@connectrpc/connect-query";
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
  getAPIKey,
  getUser,
  listAuditLogEvents,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { ListAuditLogEventsRequest } from "@/gen/tesseral/frontend/v1/frontend_pb";
import { AuditLogEvent } from "@/gen/tesseral/frontend/v1/models_pb";
import { cn } from "@/lib/utils";

import { ValueCopier } from "../core/ValueCopier";
import { TableSkeleton } from "../skeletons/TableSkeleton";
import { Badge } from "../ui/badge";
import { Button } from "../ui/button";
import { Calendar } from "../ui/calendar";
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
  listParams: ListAuditLogEventsRequest;
}) {
  const [date, setDate] = React.useState<DateRange | undefined>(undefined);
  const [eventName, setEventName] = useState<string>("");
  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});

  const stableListParams = useMemo(() => {
    const params = {
      ...listParams,
      pageToken: "",
    } as ListAuditLogEventsRequest;

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
    isLoading,
  } = useInfiniteQuery(listAuditLogEvents, stableListParams, {
    pageParamKey: "pageToken",
    getNextPageParam: (lastPage) => lastPage.nextPageToken || undefined,
  });

  const auditLogEvents =
    listAuditLogEventsResponse?.pages.flatMap((page) => page.auditLogEvents) ||
    [];
  const eventNames = Array.from(
    new Set(auditLogEvents.map((event) => event.eventName)),
  );

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
        <div className="flex flex-col lg:flex-row items-center justify-start gap-2 lg:gap-2">
          <div className="space-y-1 w-full lg:w-auto">
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant={"outline"}
                  className="justify-start text-left font-normal w-full lg:w-auto"
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

          <div className="space-y-2 w-full lg:w-auto">
            <Select value={eventName} onValueChange={setEventName}>
              <SelectTrigger className="bg-white hover:bg-muted w-full lg:w-auto">
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
                  {eventNames.map((eventName) => (
                    <SelectItem className="" key={eventName} value={eventName}>
                      {eventName}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          {((date && date.from) || (eventName && eventName.length > 0)) && (
            <div className="ml-auto w-full lg:w-auto">
              <Button
                className="w-full lg:w-auto"
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
      {isLoading ? (
        <TableSkeleton columns={4} />
      ) : (
        <>
          {auditLogEvents.length === 0 ? (
            <div className="text-sm text-center text-muted-foreground py-6">
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
        </>
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
  const { actorApiKeyId, actorUserId } = auditLogEvent;

  const { data: getApiKeyResponse } = useQuery(
    getAPIKey,
    {
      id: actorApiKeyId,
    },
    {
      enabled: !!actorApiKeyId,
      retry: false,
    },
  );
  const { data: getUserResponse } = useQuery(
    getUser,
    {
      id: actorUserId,
    },
    {
      enabled: !!actorUserId,
      retry: false,
    },
  );

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [apiKeyActor, setApiKeyActor] = useState<Record<string, any>>();
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [userActor, setUserActor] = useState<Record<string, any>>();

  useEffect(() => {
    if (getApiKeyResponse?.apiKey) {
      const apiKey = getApiKeyResponse.apiKey;
      const apiKeyActor = {
        ...apiKey,
        createTime: timestampDate(apiKey.createTime!).toISOString(),
        updateTime: timestampDate(apiKey.updateTime!).toISOString(),
      };
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      delete (apiKeyActor as any).$typeName;
      setApiKeyActor(apiKeyActor);
    }
  }, [getApiKeyResponse, setApiKeyActor]);

  useEffect(() => {
    if (getUserResponse?.user) {
      const user = getUserResponse.user;
      const userActor = {
        ...user,
        createTime: timestampDate(user.createTime!).toISOString(),
        updateTime: timestampDate(user.updateTime!).toISOString(),
      };
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      delete (userActor as any).$typeName;
      setUserActor(userActor);
    }
  }, [getUserResponse, setUserActor]);

  return (
    <>
      {apiKeyActor || userActor ? (
        <Badge variant="secondary">
          {apiKeyActor && (
            <span className="font-mono">
              {apiKeyActor.displayName || apiKeyActor.id}
            </span>
          )}
          {userActor && <span className="font-mono">{userActor.email}</span>}
        </Badge>
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
                <div className="font-mono text-xs whitespace-pre-wrap break-all p-2 rounded bg-muted border">
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
  const { actorApiKeyId, actorUserId } = event;

  const { data: getApiKeyResponse } = useQuery(
    getAPIKey,
    {
      id: actorApiKeyId,
    },
    {
      enabled: !!actorApiKeyId,
    },
  );
  const { data: getUserResponse } = useQuery(
    getUser,
    {
      id: actorUserId,
    },
    {
      enabled: !!actorUserId,
    },
  );

  return (
    <div className="space-y-4">
      {getApiKeyResponse?.apiKey && (
        <div className="space-y-1">
          <div className="font-semibold text-base">API Key</div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/organizations/${getApiKeyResponse.apiKey.organizationId}/api-keys/${getApiKeyResponse.apiKey.id}`}
          >
            {getApiKeyResponse.apiKey.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}
      {getUserResponse?.user && (
        <div className="space-y-1">
          <div className="font-semibold text-base">User</div>
          {getUserResponse.user.displayName && (
            <div className="font-medium">
              {getUserResponse.user.displayName}
            </div>
          )}
          <div
            className={cn(
              getUserResponse.user.displayName
                ? "text-muted-foreground"
                : "font-medium",
            )}
          >
            {getUserResponse.user.email}
          </div>
          <Link
            className="inline-flex items-center gap-1 text-xs font-mono px-2 py-1 rounded border text-muted-foreground hover:text-foreground bg-white"
            to={`/organization/users/${getUserResponse.user.id}`}
          >
            {getUserResponse.user.id} <ExternalLink className="h-3 w-3" />
          </Link>
        </div>
      )}

      {!getApiKeyResponse?.apiKey && !getUserResponse?.user && (
        <div className="text-muted-foreground text-sm">System</div>
      )}
    </div>
  );
}
