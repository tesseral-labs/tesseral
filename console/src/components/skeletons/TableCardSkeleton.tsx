import React from "react";

import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { TableSkeleton } from "./TableSkeleton";

export function TableCardSkeleton({
  columns = 5,
  noAction = false,
  rows = 3,
}: {
  columns?: number;
  noAction?: boolean;
  rows?: number;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="h-6 w-32 bg-gray-200 animate-pulse rounded" />
        <CardDescription className="h-4 w-md bg-gray-200 animate-pulse rounded" />
        {!noAction && (
          <CardAction>
            <div className="h-10 w-40 bg-gray-200 animate-pulse rounded-md" />
          </CardAction>
        )}
      </CardHeader>
      <CardContent>
        <TableSkeleton columns={columns} rows={rows} />
      </CardContent>
    </Card>
  );
}
