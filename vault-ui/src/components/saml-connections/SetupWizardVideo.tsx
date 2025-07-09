import clsx from "clsx";
import React, { useState } from "react";

import { Skeleton } from "@/components/ui/skeleton";

export function SetupWizardVideo({ src }: { src: string }) {
  const [loaded, setLoaded] = useState(false);

  const width = 1144;
  const height = 720;
  const aspectRatio = `${width}/${height}`;

  return (
    <div className="relative w-full" style={{ aspectRatio }}>
      {!loaded && (
        <Skeleton className="rounded-xl border shadow-md w-full h-full absolute top-0 left-0" />
      )}
      <img
        className={clsx(
          "rounded-xl border shadow-md w-full h-full object-cover",
          { hidden: !loaded },
        )}
        src={src}
        onLoad={() => setLoaded(true)}
      />
    </div>
  );
}
