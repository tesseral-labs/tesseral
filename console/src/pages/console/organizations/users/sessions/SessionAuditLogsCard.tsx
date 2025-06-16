import React from "react";

export function SessionAuditLogsCard() {
  return (
    <div className="card">
      <div className="card-header">
        <h2 className="card-title">Session Audit Logs</h2>
        <p className="card-description">
          View and manage session audit logs for users in your organization.
        </p>
      </div>
      <div className="card-content">
        <p className="text-muted-foreground text-sm">
          No session audit logs available at the moment.
        </p>
      </div>
    </div>
  );
}
