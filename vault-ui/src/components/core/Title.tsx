import { useQuery } from "@connectrpc/connect-query";
import React from "react";
import { Helmet } from "react-helmet";

import { getProject } from "@/gen/tesseral/frontend/v1/frontend-FrontendService_connectquery";

export function Title({ title }: { title?: string }) {
  const { data: getProjectResponse } = useQuery(getProject);

  return (
    <Helmet>
      {title && (
        <title>
          {title}
          {getProjectResponse?.project?.displayName
            ? `| ${getProjectResponse?.project?.displayName}`
            : ""}
        </title>
      )}
    </Helmet>
  );
}
