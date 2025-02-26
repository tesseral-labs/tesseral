import React, { forwardRef, useMemo, useState } from 'react';
import { HexColorPicker } from 'react-colorful';
import { cn } from '@/lib/utils';
import { useForwardedRef } from '@/hooks/use-forwarded-ref';
import type { ButtonProps } from '@/components/ui/button';
import { Button } from '@/components/ui/button';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Input } from '@/components/ui/input';

interface ColorPickerProps {
  value: string;
  onChange: (value: string) => void;
  onBlur?: () => void;
}

const ColorPicker = forwardRef<
  HTMLInputElement,
  Omit<ButtonProps, 'value' | 'onChange' | 'onBlur'> & ColorPickerProps
>(
  (
    { disabled, value, onChange, onBlur, name, className, ...props },
    forwardedRef,
  ) => {
    const ref = useForwardedRef(forwardedRef);
    const [open, setOpen] = useState(false);

    const parsedValue = useMemo(() => {
      return value || '#FFFFFF';
    }, [value]);

    return (
      <div className={cn('flex items-center space-x-2', className)}>
        <Popover onOpenChange={setOpen} open={open}>
          <PopoverTrigger asChild disabled={disabled} onBlur={onBlur}>
            <Button
              {...props}
              className="inline-block h-8 w-8"
              name={name}
              onClick={() => {
                setOpen(true);
              }}
              size="sm"
              style={{
                backgroundColor: parsedValue,
              }}
              variant="outline"
            >
              <div />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-full">
            <HexColorPicker color={parsedValue} onChange={onChange} />
            <Input
              maxLength={7}
              onChange={(e) => {
                onChange(e?.currentTarget?.value);
              }}
              ref={ref}
              value={parsedValue}
            />
          </PopoverContent>
        </Popover>
        <Input
          className="h-8 disabled:cursor-default disabled:bg-background disabled:border-input disabled:text-foreground disabled:opacity-100 max-w-[120px]"
          disabled
          type="text"
          value={value}
        />
      </div>
    );
  },
);
ColorPicker.displayName = 'ColorPicker';

export { ColorPicker };
