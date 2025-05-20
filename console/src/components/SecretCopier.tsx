import React, { useCallback, useEffect, useState } from 'react';
import { offset, useFloating, useTransitionStyles } from '@floating-ui/react';
import { CopyIcon } from 'lucide-react';

export const SecretCopier = ({
  placeholder,
  secret,
}: {
  placeholder: string;
  secret: string;
}) => {
  const [open, setOpen] = useState(false);
  const { refs, floatingStyles, context } = useFloating({
    open,
    onOpenChange: setOpen,
    placement: 'top',
    middleware: [offset(5)],
  });
  const { isMounted, styles } = useTransitionStyles(context, {
    duration: 150,
    initial: { opacity: 0, transform: 'translateY(0)' },
    open: { opacity: 1, transform: 'translateY(-5px)' },
  });

  useEffect(() => {
    if (open) {
      const timeoutId = setTimeout(() => {
        setOpen(false);
      }, 1000);
      return () => clearTimeout(timeoutId);
    }
  }, [open]);

  const handleCopy = useCallback(async () => {
    await navigator.clipboard.writeText(secret);
    setOpen(true);
  }, [secret, setOpen]);

  return (
    <div
      ref={refs.setReference}
      onClick={handleCopy}
      className="flex select-none cursor-pointer bg-muted font-mono text-xs border border-input rounded-md px-3 py-2 pr-12 max-w-full overflow-x-hidden whitespace-nowrap relative"
    >
      <span>{placeholder}</span>
      <div className="ml-auto flex gap-x-2 absolute right-0 top-0 p-2 bg-muted">
        <CopyIcon className="cursor-pointer text-muted-foreground hover:text-foreground h-4 w-4" />
      </div>

      {open && (
        <div ref={refs.setFloating} style={floatingStyles}>
          {isMounted && (
            <div
              style={styles}
              className="font-sans bg-black text-white px-2 py-1 text-xs rounded"
            >
              Copied!
            </div>
          )}
        </div>
      )}
    </div>
  );
};
