import { Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";

export function stepCompleted(step: Timestamp | undefined): boolean {
  if (!step) return false;

  const stepCompletionTime = timestampDate(step);

  if (stepCompletionTime.getTime() > 0) {
    return true;
  }

  return false;
}
