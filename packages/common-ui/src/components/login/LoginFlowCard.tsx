import React from "react";

import { Card } from "@/components/ui/card";

export function LoginFlowCard({ children }: { children?: React.ReactNode }) {
  return <Card className="w-full">{children}</Card>;
}
