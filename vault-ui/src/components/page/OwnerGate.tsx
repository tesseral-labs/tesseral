import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Outlet, useNavigate } from "react-router";

import { whoami } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function OwnerGate() {
  const navigate = useNavigate();
  const { data: whoamiResponse, isLoading } = useQuery(whoami);

  const user = whoamiResponse?.user;
  if (!isLoading && user && !user.owner) {
    navigate("/user");
    return null;
  }

  return <Outlet />;
}
