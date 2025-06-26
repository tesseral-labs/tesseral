import { motion } from "framer-motion";
import React, { PropsWithChildren, useEffect, useRef, useState } from "react";
import { Link } from "react-router";

import { cn } from "@/lib/utils";

export function Tabs({
  children,
  className = "",
}: PropsWithChildren & { className?: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [indicatorStyle, setIndicatorStyle] = useState({ left: 0, width: 0 });

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const activeEl = container.querySelector(
      "[data-active='true']",
    ) as HTMLElement;
    if (activeEl) {
      setIndicatorStyle({
        left: activeEl.offsetLeft,
        width: activeEl.offsetWidth,
      });
    }
  }, [children]);

  return (
    <div
      className={cn(
        "relative inline-block bg-gradient-to-br from-gray-100/50 to-gray-200/50 rounded-md p-1",
        className,
      )}
    >
      <div ref={containerRef} className="relative flex space-x-2">
        <motion.div
          layout
          transition={{ type: "spring", stiffness: 400, damping: 30 }}
          className="absolute top-0 bottom-0 bg-white rounded-sm shadow-sm"
          style={{
            left: indicatorStyle.left,
            width: indicatorStyle.width,
          }}
        />
        {children}
      </div>
    </div>
  );
}

export function Tab({
  active = false,
  children,
}: PropsWithChildren<{ active?: boolean }>) {
  return (
    <div
      className={cn(
        "relative z-1 inline-block px-4 py-2 text-sm rounded-sm cursor-pointer transition-colors ",
        active
          ? "bg-white font-medium shadow-sm"
          : "text-muted-foreground hover:shadow-sm hover:bg-white/50 hover:text-foreground/80",
      )}
      data-active={active || undefined}
    >
      {children}
    </div>
  );
}

export function TabLink({
  active = false,
  children,
  to,
}: PropsWithChildren<{
  active?: boolean;
  to: string;
}>) {
  return (
    <Link to={to}>
      <Tab active={active}>{children}</Tab>
    </Link>
  );
}
