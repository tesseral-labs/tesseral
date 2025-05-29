// src/components/audit-log-viewer.tsx
import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import { useQuery } from "@connectrpc/connect-query";
import { format } from "date-fns";
import {
  ArrowLeft,
  ArrowRight,
  CalendarIcon,
  ChevronDown,
  ChevronRight,
  FilterX,
  Search,
} from "lucide-react";
import React, { useCallback, useState } from "react";
import { DateRange } from "react-day-picker";

import { MultiSelect } from "@/components/MultiSelect";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Input } from "@/components/ui/input";
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
import { AuditLogEvent } from "@/gen/tesseral/common/v1/common_pb";
import {
  getUser,
  listAuditLogEvents,
} from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import { getAPIKey } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";
import {
  ListAuditLogEventsRequest,
  ListAuditLogEventsRequest_Filter,
} from "@/gen/tesseral/frontend/v1/frontend_pb";

const PAGE_SIZE = 10;

// --- Filter Bar Component ---
interface FilterBarProps {
  onApply: (filter: ListAuditLogEventsRequest_Filter) => void;
  isLoading: boolean;
}

function makeFilter(
  params: Omit<ListAuditLogEventsRequest_Filter, "$typeName" | "eventName"> & {
    eventName?: string[];
  },
): ListAuditLogEventsRequest_Filter {
  return {
    $typeName: "tesseral.frontend.v1.ListAuditLogEventsRequest.Filter",
    ...params,
    eventName: params.eventName ?? [],
  };
}

function FilterBar({ onApply, isLoading }: FilterBarProps) {
  const [date, setDate] = React.useState<DateRange | undefined>(undefined);
  const [eventNames, setEventNames] = useState<string[]>([]);
  const [userId, setUserId] = useState("");

  function handleApply() {
    const filter: ListAuditLogEventsRequest_Filter = makeFilter({});
    if (date?.from) {
      filter.startTime = timestampFromDate(date.from);
    }
    if (date?.to) {
      // Set endTime to the end of the selected day (23:59:59.999)
      const end = new Date(date.to);
      end.setHours(23, 59, 59, 999);
      filter.endTime = timestampFromDate(end);
    }
    if (eventNames.length > 0) filter.eventName = eventNames;
    if (userId) filter.userId = userId;

    onApply(filter);
  }

  function handleReset() {
    setDate(undefined);
    setEventNames([]);
    setUserId("");
    onApply(makeFilter({}));
  }

  const hasFilters = date || eventNames.length > 0 || userId;

  return (
    <div className="p-4 border-b bg-card">
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

        {/* Event Name Selector */}
        <MultiSelect
          selected={eventNames}
          onChange={setEventNames}
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
          <Button onClick={handleApply} disabled={isLoading}>
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

// --- Main Viewer Component ---
export function AuditLogEventsPage() {
  const [request, setRequest] = useState<ListAuditLogEventsRequest>({
    $typeName: "tesseral.frontend.v1.ListAuditLogEventsRequest",
    pageSize: PAGE_SIZE,
    pageToken: "",
  });
  const [pageTokens, setPageTokens] = useState<string[]>([""]);
  const [currentPageIndex, setCurrentPageIndex] = useState(0);

  // Track expanded rows by event ID
  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});

  function toggleRow(eventId: string) {
    setExpandedRows((prev) => ({
      ...prev,
      [eventId]: !prev[eventId],
    }));
  }

  const handleApplyFilters = useCallback(
    (filter: ListAuditLogEventsRequest_Filter) => {
      setRequest({
        $typeName: "tesseral.frontend.v1.ListAuditLogEventsRequest",
        pageSize: PAGE_SIZE,
        filter: Object.keys(filter).length > 0 ? filter : undefined,
        pageToken: "",
      });
      setPageTokens([""]); // Reset pagination on filter change
      setCurrentPageIndex(0);
      setExpandedRows({}); // Reset expanded rows on filter change
    },
    [],
  );

  const currentRequest: ListAuditLogEventsRequest = {
    ...request,
    pageToken: pageTokens[currentPageIndex] ?? "",
  };

  const { data, isLoading, isError, error, isFetching } = useQuery(
    listAuditLogEvents,
    currentRequest,
  );

  function handleNextPage() {
    if (data?.nextPageToken) {
      // If we've been here before, just move forward
      if (currentPageIndex < pageTokens.length - 1) {
        setCurrentPageIndex(currentPageIndex + 1);
      } else {
        // Otherwise, add the new token and move
        setPageTokens([...pageTokens, data.nextPageToken]);
        setCurrentPageIndex(currentPageIndex + 1);
      }
      setExpandedRows({}); // Reset expanded rows on next page
    }
  }

  function handlePrevPage() {
    if (currentPageIndex > 0) {
      setCurrentPageIndex(currentPageIndex - 1);
    }
  }

  function renderActor(event: AuditLogEvent) {
    if (event.userId) {
      return event.userId;
    }
    if (event.apiKeyId) {
      return event.apiKeyId;
    }
    return <span className="text-muted-foreground">System</span>;
  }

  return (
    <TooltipProvider>
      <div className="border rounded-lg shadow-sm">
        <FilterBar
          onApply={handleApplyFilters}
          isLoading={isLoading || isFetching}
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
              {isLoading && !data && (
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
              {!isLoading && !isError && data?.auditLogEvents.length === 0 && (
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
                data?.auditLogEvents.map((event) => (
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
                            {event.userId && <p>User ID: {event.userId}</p>}
                            {event.sessionId && (
                              <p>Session ID: {event.sessionId}</p>
                            )}
                            {event.apiKeyId && (
                              <p>API Key ID: {event.apiKeyId}</p>
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
        {/* Pagination Controls */}
        <div className="flex items-center justify-end space-x-2 p-4 border-t">
          <span className="text-sm text-muted-foreground">
            Page {currentPageIndex + 1}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={handlePrevPage}
            disabled={currentPageIndex === 0 || isLoading || isFetching}
          >
            <ArrowLeft className="h-4 w-4 mr-1" /> Previous
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleNextPage}
            disabled={!data?.nextPageToken || isLoading || isFetching}
          >
            Next <ArrowRight className="h-4 w-4 ml-1" />
          </Button>
        </div>
      </div>
    </TooltipProvider>
  );
}

function AuditLogEventDetails({ event }: { event: AuditLogEvent }) {
  let actorDetails: React.ReactNode = null;
  if (event.apiKeyId) {
    actorDetails = <AuditLogEventApiKeyDetails apiKeyId={event.apiKeyId} />;
  } else if (event.userId) {
    actorDetails = <AuditLogEventUserDetails userId={event.userId} />;
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
