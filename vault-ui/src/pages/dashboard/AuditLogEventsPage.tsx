// src/components/audit-log-viewer.tsx
import React, { useState, useCallback } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Skeleton } from "@/components/ui/skeleton";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { CalendarIcon, ArrowLeft, ArrowRight, X, FilterX, ArrowUpDown, Search, ShieldIcon } from 'lucide-react';
import { format } from 'date-fns';
import { DateRange } from 'react-day-picker';
import { timestampDate, timestampFromDate } from '@bufbuild/protobuf/wkt';
import { ListAuditLogEventsRequest, ListAuditLogEventsRequest_Filter, ListAuditLogEventsRequest_FilterSchema } from '@/gen/tesseral/frontend/v1/frontend_pb';
import { MultiSelect } from '@/components/MultiSelect';
import { useQuery } from '@connectrpc/connect-query';
import { listAuditLogEvents } from '@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery';
import { AuditLogEvent } from '@/gen/tesseral/common/v1/common_pb';

const PAGE_SIZE = 20;

// --- Mock Event Names (Replace with actual fetch if dynamic) ---
const availableEventNames = [
    { value: 'user.login', label: 'User Login' },
    { value: 'user.logout', label: 'User Logout' },
    { value: 'file.create', label: 'File Created' },
    { value: 'file.delete', label: 'File Deleted' },
    { value: 'file.update', label: 'File Updated' },
    { value: 'settings.update', label: 'Settings Updated' },
];

// --- Filter Bar Component ---
interface FilterBarProps {
    onApply: (filter: ListAuditLogEventsRequest_Filter, orderBy: string) => void;
    isLoading: boolean;
}

function makeFilter(params: Omit<ListAuditLogEventsRequest_Filter, '$typeName' | 'eventName'> & {
    eventName?: string[];
}): ListAuditLogEventsRequest_Filter {
    return {
        $typeName: 'tesseral.frontend.v1.ListAuditLogEventsRequest.Filter',
        ...params,
        eventName: params.eventName ?? [],
    };
}

const FilterBar: React.FC<FilterBarProps> = ({ onApply, isLoading }) => {
    const [date, setDate] = React.useState<DateRange | undefined>(undefined);
    const [eventNames, setEventNames] = useState<string[]>([]);
    const [userId, setUserId] = useState('');
    const [orderBy, setOrderBy] = useState('id desc'); // Default sort

    const handleApply = () => {
        const filter: ListAuditLogEventsRequest_Filter = makeFilter({});
        if (date?.from) {
            filter.startTime = timestampFromDate(date.from);
            // NOTE: Protobuf message only has start_time. If you need end_time,
            // you'd need to adjust your API or handle filtering client-side
            // (not ideal for pagination). We'll only use start_time here.
        }
        if (eventNames.length > 0) filter.eventName = eventNames;
        if (userId) filter.userId = userId;

        onApply(filter, orderBy);
    };

    const handleReset = () => {
        setDate(undefined);
        setEventNames([]);
        setUserId('');
        setOrderBy('id desc');
        onApply(makeFilter({}), 'id desc');
    };

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
                                        {format(date.from, "LLL dd, y")} - {format(date.to, "LLL dd, y")}
                                    </>
                                ) : (
                                    format(date.from, "LLL dd, y")
                                )
                            ) : (
                                <span>Pick a start date</span>
                            )}
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                            initialFocus
                            mode="range" // Use range, but only send 'from' based on proto
                            selected={date}
                            onSelect={setDate}
                            numberOfMonths={1}
                        />
                    </PopoverContent>
                </Popover>

                {/* Event Name Selector */}
                <MultiSelect
                    options={availableEventNames}
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
                       <Button variant="outline" onClick={handleReset} disabled={isLoading}>
                           <FilterX className="mr-2 h-4 w-4" /> Reset
                       </Button>
                    )}
                </div>
            </div>
             <p className="text-xs text-muted-foreground mt-2">
                Sorting by: {orderBy}. (Click headers to change - requires API support)
            </p>
        </div>
    );
};

// --- Main Viewer Component ---
export function AuditLogEventsPage() {
    const [request, setRequest] = useState<ListAuditLogEventsRequest>({
        $typeName: 'tesseral.frontend.v1.ListAuditLogEventsRequest',
        pageSize: PAGE_SIZE,
        pageToken: '',
        orderBy: 'id desc',
    });
    const [pageTokens, setPageTokens] = useState<string[]>(['']); // Start with an empty token for page 1
    const [currentPageIndex, setCurrentPageIndex] = useState(0);

    const handleApplyFilters = useCallback((filter: ListAuditLogEventsRequest_Filter, orderBy: string) => {
        setRequest({
            $typeName: 'tesseral.frontend.v1.ListAuditLogEventsRequest',
            pageSize: PAGE_SIZE,
            filter: Object.keys(filter).length > 0 ? filter : undefined,
            orderBy: orderBy,
            pageToken: '', // TODO
        });
        setPageTokens(['']); // Reset pagination on filter change
        setCurrentPageIndex(0);
    }, []);

    const currentRequest = {
        ...request,
        pageToken: pageTokens[currentPageIndex] || undefined,
    };

    const { data, isLoading, isError, error, isFetching } = useQuery(
        listAuditLogEvents
    );

    const handleNextPage = () => {
        if (data?.nextPageToken) {
            // If we've been here before, just move forward
            if (currentPageIndex < pageTokens.length - 1) {
                 setCurrentPageIndex(currentPageIndex + 1);
            } else {
            // Otherwise, add the new token and move
                 setPageTokens([...pageTokens, data.nextPageToken]);
                 setCurrentPageIndex(currentPageIndex + 1);
            }
        }
    };

    const handlePrevPage = () => {
        if (currentPageIndex > 0) {
            setCurrentPageIndex(currentPageIndex - 1);
        }
    };

    const handleSort = (field: 'id' | 'event_name' | 'user_id') => {
       const currentOrderBy = request.orderBy || 'id desc';
       let newOrderBy: string;

       if (currentOrderBy.startsWith(field)) {
           newOrderBy = currentOrderBy.endsWith('desc') ? `${field} asc` : `${field} desc`;
       } else {
           newOrderBy = `${field} desc`; // Default to desc when changing field
       }
       // IMPORTANT: Ensure your backend supports these order_by strings!
       // The SQL only explicitly indexes `id desc`. Other sorts might be slow
       // or require additional indexes or backend logic.
       console.warn(`Sorting by ${newOrderBy}. Ensure your backend supports this.`);
       handleApplyFilters(request.filter || makeFilter({}), newOrderBy);
    };


    const renderActor = (event: AuditLogEvent) => {
        if (event.userId) return `User: ${event.userId.substring(0, 8)}...`;
        if (event.sessionId) return `Session: ${event.sessionId.substring(0, 8)}...`;
        if (event.apiKeyId) return `API Key: ${event.apiKeyId.substring(0, 8)}...`;
        return <span className="text-muted-foreground">System</span>;
    };

    const getSortIndicator = (field: string) => {
        const currentOrderBy = request.orderBy || 'id desc';
        if (currentOrderBy.startsWith(field)) {
             return currentOrderBy.endsWith('desc') ? ' ▼' : ' ▲';
        }
        return <ArrowUpDown className="ml-2 h-3 w-3 inline-block opacity-30 group-hover:opacity-100" />;
    }

    return (
        <TooltipProvider>
            <div className="border rounded-lg shadow-sm">
                <FilterBar onApply={handleApplyFilters} isLoading={isLoading || isFetching} />

                <div className="p-4">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-[50px]"></TableHead>
                                <TableHead
                                    className="cursor-pointer group"
                                    onClick={() => handleSort('event_name')}
                                >
                                    Event {getSortIndicator('event_name')}
                                </TableHead>
                                <TableHead>Actor</TableHead>
                                <TableHead
                                     className="cursor-pointer group"
                                     onClick={() => handleSort('id')} // Sorting by ID = Sorting by time
                                >
                                    Time {getSortIndicator('id')}
                                </TableHead>
                                <TableHead>Details</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {isLoading && !data && (
                                <>
                                    {Array.from({ length: 10 }).map((_, i) => (
                                        <TableRow key={`skel-${i}`}>
                                            <TableCell><Skeleton className="h-4 w-4 rounded-full" /></TableCell>
                                            <TableCell><Skeleton className="h-4 w-[150px]" /></TableCell>
                                            <TableCell><Skeleton className="h-4 w-[100px]" /></TableCell>
                                            <TableCell><Skeleton className="h-4 w-[200px]" /></TableCell>
                                            <TableCell><Skeleton className="h-4 w-[50px]" /></TableCell>
                                        </TableRow>
                                    ))}
                                </>
                            )}
                            {isError && (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center text-destructive">
                                        Failed to load audit logs: { (error as Error)?.message || 'Unknown error' }
                                    </TableCell>
                                </TableRow>
                            )}
                            {!isLoading && !isError && data?.auditLogEvents.length === 0 && (
                                 <TableRow>
                                    <TableCell colSpan={5} className="text-center text-muted-foreground">
                                        No audit log events found matching your criteria.
                                    </TableCell>
                                </TableRow>
                            )}
                            {!isLoading && data?.auditLogEvents.map((event) => (
                                <TableRow key={event.id}>
                                    <TableCell>
                                         <Tooltip>
                                            <TooltipTrigger>
                                                <AuditLogIcon eventName={event.eventName} />
                                            </TooltipTrigger>
                                            <TooltipContent>
                                                <p>{event.eventName}</p>
                                            </TooltipContent>
                                        </Tooltip>
                                    </TableCell>
                                    <TableCell className="font-medium">{event.eventName}</TableCell>
                                    <TableCell>
                                        <Tooltip>
                                            <TooltipTrigger>
                                                <span>{renderActor(event)}</span>
                                            </TooltipTrigger>
                                            <TooltipContent>
                                                {event.userId && <p>User ID: {event.userId}</p>}
                                                {event.sessionId && <p>Session ID: {event.sessionId}</p>}
                                                {event.apiKeyId && <p>API Key ID: {event.apiKeyId}</p>}
                                            </TooltipContent>
                                        </Tooltip>
                                    </TableCell>
                                    <TableCell>
                                        {format(timestampDate(event.eventTime!), 'PPpp')}
                                    </TableCell>
                                     <TableCell>
                                         <Popover>
                                            <PopoverTrigger asChild>
                                                <Button variant="ghost" size="sm">View</Button>
                                            </PopoverTrigger>
                                            <PopoverContent className="w-80">
                                                <pre className="text-xs overflow-auto max-h-60">
                                                    {JSON.stringify(event.eventDetails, null, 2)}
                                                </pre>
                                            </PopoverContent>
                                        </Popover>
                                    </TableCell>
                                </TableRow>
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
                        <ArrowLeft className="h-4 w-4 mr-1"/> Previous
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleNextPage}
                        disabled={!data?.nextPageToken || isLoading || isFetching}
                    >
                        Next <ArrowRight className="h-4 w-4 ml-1"/>
                    </Button>
                </div>
            </div>
        </TooltipProvider>
    );
};

interface AuditLogIconProps {
  eventName: string;
  className?: string;
}

const AuditLogIcon: React.FC<AuditLogIconProps> = ({ eventName, className }) => {
    // Try to find a specific match, then a prefix match, then default
    const IconComponent = ShieldIcon;
    return <IconComponent className={className || "h-4 w-4 text-muted-foreground"} />;
};