import { VariantProps, cva } from "class-variance-authority";
import { LoaderCircle } from "lucide-react";
import React, { SVGAttributes, forwardRef } from "react";

import { cn } from "@/lib/utils";

const loaderVariants = cva("animate-spin", {
  variants: {
    variant: {
      default: "text-foreground",
      primary: "text-primary",
    },
    size: {
      default: "w-6 h-6",
      sm: "w-4 h-4",
      lg: "w-8 h-8",
    },
  },
});

export interface LoaderProps
  extends SVGAttributes<SVGElement>,
    VariantProps<typeof loaderVariants> {
  asChild?: boolean;
}

const Loader = forwardRef<SVGElement, LoaderProps>(
  ({ className, size, variant }) => {
    return (
      <LoaderCircle
        className={cn(loaderVariants({ variant, size, className }))}
      />
    );
  },
);
Loader.displayName = "Loader";

export default Loader;
