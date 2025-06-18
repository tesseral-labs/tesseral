import { CircleXIcon } from "lucide-react";
import React from "react";

import { StatusIndicator } from "@/components/core/StatusIndicator";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { VaultDomainSettingsDNSRecord } from "@/gen/tesseral/backend/v1/models_pb";

export function DNSRecordRows({
  record,
}: {
  record: VaultDomainSettingsDNSRecord;
}) {
  const noValue = (record.actualValues ?? []).length === 0;
  const tooManyValues = record.actualValues?.length > 1;
  const incorrectValue =
    record.actualValues?.length === 1 &&
    record.actualValues[0] !== record.wantValue;

  return (
    <>
      <TableRow>
        <TableCell>
          {record.correct && (
            <StatusIndicator variant="success">Configured</StatusIndicator>
          )}
          {noValue && (
            <StatusIndicator variant="pending">No record</StatusIndicator>
          )}
          {(tooManyValues || incorrectValue) && (
            <StatusIndicator variant="error">Misconfigured</StatusIndicator>
          )}
        </TableCell>
        <TableCell>{record.type}</TableCell>
        <TableCell>{record.name}</TableCell>
        <TableCell>{record.wantValue}</TableCell>
      </TableRow>

      {incorrectValue && (
        <TableRow className="bg-red-50/50 hover:bg-red-50/50">
          <TableCell colSpan={4}>
            <Alert variant="destructive" className="bg-white">
              <CircleXIcon className="w-5 h-5 text-red-500" />
              <AlertTitle>
                <span className="font-mono">{record.name}</span> is
                misconfigured
              </AlertTitle>
              <AlertDescription>
                <p className="mt-2">This record has the wrong value.</p>

                <Table>
                  <TableBody>
                    <TableRow className="border-destructive/25 hover:bg-white">
                      <TableCell>Expected</TableCell>
                      <TableCell className="font-mono">
                        {record.wantValue}
                      </TableCell>
                    </TableRow>
                    <TableRow className="border-destructive/25 hover:bg-white">
                      <TableCell>Actual</TableCell>
                      <TableCell className="font-mono">
                        {record.actualValues[0]}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
                <p className="mt-2">
                  It will take at least {record.actualTtlSeconds} seconds for
                  any change you make here to propagate, because that's the
                  time-to-live (TTL) you configured on this incorrect record.
                </p>
              </AlertDescription>
            </Alert>
          </TableCell>
        </TableRow>
      )}

      {tooManyValues && (
        <TableRow className="bg-red-50/50 hover:bg-red-50/50">
          <TableCell colSpan={4}>
            <Alert variant="destructive" className="bg-white">
              <CircleXIcon className="w-5 h-5 text-red-500" />
              <AlertTitle>
                <span className="font-mono">{record.name}</span> is
                misconfigured
              </AlertTitle>
              <AlertDescription>
                <p className="mt-2">
                  This record has too many values. Delete the following records:
                </p>

                <Table>
                  <TableHeader className="border-destructive/25 border-b">
                    <TableRow className="hover:bg-white">
                      <TableHead className="text-destructive">Type</TableHead>
                      <TableHead className="text-destructive">Name</TableHead>
                      <TableHead className="text-destructive">Value</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {record.actualValues
                      ?.filter((v) => v !== record.wantValue)
                      .map((v, i) => (
                        <TableRow
                          key={i}
                          className="border-destructive/25 hover:bg-white"
                        >
                          <TableCell>{record.type}</TableCell>
                          <TableCell>{record.name}</TableCell>
                          <TableCell>{v}</TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
                <p className="mt-2">
                  It will take at least {record.actualTtlSeconds} seconds for
                  any change you make here to propagate, because that's the
                  time-to-live (TTL) you configured on these records.
                </p>
              </AlertDescription>
            </Alert>
          </TableCell>
        </TableRow>
      )}
    </>
  );
}
