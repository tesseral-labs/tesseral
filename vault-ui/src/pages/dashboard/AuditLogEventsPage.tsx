import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import { useInfiniteQuery, useQuery } from "@connectrpc/connect-query";
import { format } from "date-fns";
import {
  CalendarIcon,
  ChevronDown,
  ChevronRight,
  FilterX,
  Search,
} from "lucide-react";
import React, {
  Dispatch,
  SetStateAction,
  useEffect,
  useMemo,
  useState,
} from "react";
import { DateRange } from "react-day-picker";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Input } from "@/components/ui/input";
import Loader from "@/components/ui/loader";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  getUser,
  listAuditLogEvents,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { getAPIKey } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { ListAuditLogEventsRequest } from "@/gen/tesseral/frontend/v1/frontend_pb";
import { AuditLogEvent } from "@/gen/tesseral/frontend/v1/models_pb";

interface FilterBarProps {
  setParams: Dispatch<SetStateAction<Partial<ListAuditLogEventsRequest>>>;
  isLoading: boolean;
}

function FilterBar({ setParams, isLoading }: FilterBarProps) {
  const [date, setDate] = React.useState<DateRange | undefined>(undefined);
  const [eventName, setEventName] = useState<string>("");
  const [userId, setUserId] = useState("");

  function handleApply() {
    const filter: Partial<ListAuditLogEventsRequest> = {};
    if (date?.from) {
      filter.filterStartTime = timestampFromDate(date.from);
    }

    if (date?.to) {
      // Set endTime to the end of the selected day (23:59:59.999)
      const end = new Date(date.to);
      end.setHours(23, 59, 59, 999);
      filter.filterEndTime = timestampFromDate(end);
    }

    if (eventName.length > 0) {
      filter.filterEventName = eventName;
    }

    if (userId) filter.filterUserId = userId;

    setParams(filter);
  }

  function handleReset() {
    setDate(undefined);
    setEventName("");
    setUserId("");
    setParams({}); // Reset all filters
  }

  const hasFilters = useMemo(
    () => date || eventName.length > 0 || userId,
    [date, eventName, userId],
  );

  useEffect(() => {
    // Reset filters if no filters are applied
    if (!hasFilters && setParams) {
      setParams({});
    }
  }, [hasFilters, setParams]);

  return (
    <div className="p-4">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Date Picker */}
        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant={"outline"}
              className="w-full justify-start text-left font-normal"
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
                <span>Pick a date range</span>
              )}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0" align="start">
            <Calendar
              initialFocus
              mode="range"
              selected={date}
              onSelect={setDate}
              numberOfMonths={1}
            />
          </PopoverContent>
        </Popover>

        {/* Event Name Input */}
        <Input
          value={eventName}
          onChange={(e) => setEventName(e.target.value)}
          placeholder="Filter by event name..."
          className="w-full"
        />

        {/* User ID Input */}
        <Input
          placeholder="Filter by User ID..."
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          className="w-full"
        />

        {/* Actions */}
        <div className="flex gap-2 justify-end">
          <Button onClick={handleApply} disabled={isLoading || !hasFilters}>
            <Search className="mr-2 h-4 w-4" /> Apply Filters
          </Button>
          {hasFilters && (
            <Button
              variant="outline"
              onClick={handleReset}
              disabled={isLoading}
            >
              <FilterX className="mr-2 h-4 w-4" /> Reset
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}

export function AuditLogEventsPage() {
  const [listAuditLogEventsParams, setListAuditLogEventsParams] = useState<
    Partial<ListAuditLogEventsRequest>
  >({});

  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});
  const stableParams = useMemo(
    () => ({ ...listAuditLogEventsParams }),
    [listAuditLogEventsParams],
  );

  const {
    data: listAuditLogEventsResponses,
    error,
    fetchNextPage,
    hasNextPage,
    isError,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery(
    listAuditLogEvents,
    {
      ...stableParams,
      pageToken: "",
    } as ListAuditLogEventsRequest,
    {
      pageParamKey: "pageToken",
      getNextPageParam: (page) => page.nextPageToken || undefined,
    },
  );

  const auditLogEvents = listAuditLogEventsResponses?.pages?.flatMap(
    (page) => page.auditLogEvents,
  );

  function toggleRow(eventId: string) {
    setExpandedRows((prev) => ({
      ...prev,
      [eventId]: !prev[eventId],
    }));
  }

  function renderActor(event: AuditLogEvent) {
    if (event.actorUserId) {
      return event.actorUserId;
    }
    if (event.actorApiKeyId) {
      return event.actorApiKeyId;
    }
    return <span className="text-muted-foreground">System</span>;
  }

  return (
    <TooltipProvider>
      <div className="border rounded-lg shadow-sm">
        <FilterBar
          setParams={setListAuditLogEventsParams}
          isLoading={isLoading}
        />
        <div className="p-4">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[40px]"></TableHead>
                <TableHead className="group">Event</TableHead>
                <TableHead>Actor</TableHead>
                <TableHead className="group">Time ▼</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading && !auditLogEvents && (
                <>
                  {Array.from({ length: 10 }).map((_, i) => (
                    <TableRow key={`skel-${i}`}>
                      <TableCell>
                        <Skeleton className="h-4 w-4 rounded-full" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-[150px]" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-[100px]" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-[200px]" />
                      </TableCell>
                    </TableRow>
                  ))}
                </>
              )}
              {isError && (
                <TableRow>
                  <TableCell
                    colSpan={4}
                    className="text-center text-destructive"
                  >
                    Failed to load audit logs:{" "}
                    {(error as Error)?.message ?? "Unknown error"}
                  </TableCell>
                </TableRow>
              )}
              {!isLoading && !isError && !auditLogEvents?.length && (
                <TableRow>
                  <TableCell
                    colSpan={4}
                    className="text-center text-muted-foreground"
                  >
                    No audit log events found matching your criteria.
                  </TableCell>
                </TableRow>
              )}
              {!isLoading &&
                auditLogEvents?.map((event) => (
                  <React.Fragment key={event.id}>
                    <TableRow
                      className="cursor-pointer"
                      onClick={() => toggleRow(event.id)}
                      data-testid={`audit-log-row-${event.id}`}
                    >
                      <TableCell className="align-middle">
                        {expandedRows[event.id] ? (
                          <ChevronDown className="h-4 w-4" />
                        ) : (
                          <ChevronRight className="h-4 w-4" />
                        )}
                      </TableCell>
                      <TableCell className="font-medium">
                        {event.eventName}
                      </TableCell>
                      <TableCell>
                        <Tooltip>
                          <TooltipTrigger>
                            <span>{renderActor(event)}</span>
                          </TooltipTrigger>
                          <TooltipContent>
                            {event.actorUserId && (
                              <p>User ID: {event.actorUserId}</p>
                            )}
                            {event.actorSessionId && (
                              <p>Session ID: {event.actorSessionId}</p>
                            )}
                            {event.actorApiKeyId && (
                              <p>API Key ID: {event.actorApiKeyId}</p>
                            )}
                          </TooltipContent>
                        </Tooltip>
                      </TableCell>
                      <TableCell>
                        <Tooltip>
                          <TooltipTrigger>
                            <span>
                              {format(timestampDate(event.eventTime!), "PPpp")}
                            </span>
                          </TooltipTrigger>
                          <TooltipContent>
                            <pre>
                              {timestampDate(event.eventTime!).toISOString()}
                            </pre>
                          </TooltipContent>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                    {expandedRows[event.id] && (
                      <AuditLogEventDetails event={event} />
                    )}
                  </React.Fragment>
                ))}
            </TableBody>
          </Table>
        </div>

        {hasNextPage && (
          <div className="flex justify-center p-4">
            <Button
              variant="outline"
              onClick={() => fetchNextPage()}
              disabled={isFetchingNextPage || isLoading}
            >
              {isFetchingNextPage ? (
                <>
                  <Loader />
                  Loading...
                </>
              ) : (
                <>Load More</>
              )}
            </Button>
          </div>
        )}
      </div>
    </TooltipProvider>
  );
}

function AuditLogEventDetails({ event }: { event: AuditLogEvent }) {
  let actorDetails: React.ReactNode = null;
  if (event.actorApiKeyId) {
    actorDetails = (
      <AuditLogEventApiKeyDetails apiKeyId={event.actorApiKeyId} />
    );
  } else if (event.actorUserId) {
    actorDetails = <AuditLogEventUserDetails userId={event.actorUserId} />;
  } else {
    actorDetails = <div className="text-muted-foreground">System</div>;
  }

  return (
    <TableRow className="bg-muted/40">
      <TableCell colSpan={4} className="p-4">
        <div className="flex flex-col md:flex-row gap-4">
          {/* Left: Actor details */}
          <div className="w-full md:w-1/2 border-r pr-4">
            <div className="font-semibold mb-2">Actor Details</div>
            {actorDetails}
          </div>
          {/* Right: Event details */}
          <div className="w-full md:w-1/2">
            <div className="font-semibold mb-2">Event Details</div>
            <div className="space-y-2">
              <div>
                <div className="text-sm font-medium">ID</div>
                <div className="text-sm">{event.id}</div>
              </div>
              <div>
                <div className="text-sm font-medium">JSON</div>
                <div className="font-mono text-xs whitespace-pre-wrap break-all">
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

function AuditLogEventUserDetails({ userId }: { userId: string }) {
  const { data, isLoading, isError, error } = useQuery(getUser, { id: userId });

  if (isLoading) {
    return <Skeleton className="h-4 w-[200px]" />;
  }
  if (isError) {
    return (
      <span className="text-destructive">
        Failed to load user details:{" "}
        {(error as Error)?.message ?? "Unknown error"}
      </span>
    );
  }

  const user = data?.user;
  if (!user) {
    return <span className="text-muted-foreground">User not found</span>;
  }

  return (
    <div className="space-y-2">
      <div>
        <div className="text-sm font-medium">ID</div>
        <div className="text-sm">{user.id}</div>
      </div>
      {user.displayName && (
        <div>
          <div className="text-sm font-medium">Display Name</div>
          <div className="text-sm">{user.displayName}</div>
        </div>
      )}
      {user.email && (
        <div>
          <div className="text-sm font-medium">Email</div>
          <div className="text-sm">{user.email}</div>
        </div>
      )}
    </div>
  );
}

function AuditLogEventApiKeyDetails({ apiKeyId }: { apiKeyId: string }) {
  const { data, isLoading, isError, error } = useQuery(getAPIKey, {
    id: apiKeyId,
  });

  if (isLoading) {
    return <Skeleton className="h-4 w-[200px]" />;
  }
  if (isError) {
    return (
      <span className="text-destructive">
        Failed to load API key details:{" "}
        {(error as Error)?.message ?? "Unknown error"}
      </span>
    );
  }
  const apiKey = data?.apiKey;
  if (!apiKey) {
    return <span className="text-muted-foreground">API Key not found</span>;
  }

  return (
    <div className="space-y-2">
      <div>
        <div className="text-sm font-medium">ID</div>
        <div className="text-sm font-mono">{apiKey.id}</div>
      </div>
      {apiKey.displayName && (
        <div>
          <div className="text-sm font-medium">Display Name</div>
          <div className="text-sm font-mono">{apiKey.displayName}</div>
        </div>
      )}
    </div>
  );
}
