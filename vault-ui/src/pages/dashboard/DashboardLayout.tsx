import React from "react";
import { Outlet } from "react-router";

export function DashboardLayout() {
  return (
    <main className="bg-body w-full">
      <div className="mx-auto max-w-7xl sm:px-6 lg:px-8 py-8">
        <Outlet />
      </div>
    </main>
  );
}
