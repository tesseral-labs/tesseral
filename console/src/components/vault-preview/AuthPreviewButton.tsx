import { VariantProps, cva } from "class-variance-authority";
import React, { forwardRef } from "react";

import { cn } from "@/lib/utils";

const authPreviewButtonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground",
        destructive: "bg-destructive text-destructive-foreground",
        outline: "border border-input bg-background text-foreground",
        secondary: "bg-secondary text-secondary-foreground",
        ghost: "",
        link: "text-primary underline-offset-4",
      },
      size: {
        default: "h-9 rounded-md px-3",
        sm: "h-9 rounded-md px-3",
        md: "h-10 px-4 py-2",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export const AuthPreviewButton = forwardRef<
  HTMLDivElement,
  React.ButtonHTMLAttributes<HTMLButtonElement> &
    VariantProps<typeof authPreviewButtonVariants>
>(({ children, className, variant }, ref) => {
  return (
    <div
      className={cn(
        "cursor-default",
        authPreviewButtonVariants({ variant }),
        className,
      )}
      ref={ref}
    >
      {children}
    </div>
  );
});
AuthPreviewButton.displayName = "AuthPreviewButton";
