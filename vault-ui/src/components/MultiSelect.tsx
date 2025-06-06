// src/components/ui/multi-select.tsx
import * as React from "react";
import { X } from "lucide-react";

import { Badge } from "@/components/ui/badge";

interface MultiSelectProps {
  selected: string[];
  onChange: (selected: string[]) => void;
  placeholder?: string;
  className?: string;
}

export function MultiSelect({
  selected,
  onChange,
  placeholder = "Select...",
  className,
}: MultiSelectProps) {
  const inputRef = React.useRef<HTMLInputElement>(null);
  const [inputValue, setInputValue] = React.useState("");

  const handleUnselect = React.useCallback(
    (value: string) => {
      onChange(selected.filter((s) => s !== value));
    },
    [onChange, selected],
  );

  const handleKeyDown = React.useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if ((e.key === "Enter" || e.key === ",") && inputValue.trim() !== "") {
        e.preventDefault();
        const value = inputValue.trim();
        if (!selected.includes(value)) {
          onChange([...selected, value]);
        }
        setInputValue("");
      } else if (
        (e.key === "Delete" || e.key === "Backspace") &&
        inputValue === "" &&
        selected.length > 0
      ) {
        handleUnselect(selected[selected.length - 1]);
      } else if (e.key === "Escape") {
        inputRef.current?.blur();
      }
    },
    [inputValue, selected, onChange, handleUnselect],
  );

  return (
    <div className={`overflow-visible bg-transparent ${className}`}>
      <div className="group border border-input px-3 py-2 text-sm ring-offset-background rounded-md focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2">
        <div className="flex gap-1 flex-wrap">
          {selected.map((value) => (
            <Badge key={value} variant="secondary">
              {value}
              <button
                className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    handleUnselect(value);
                  }
                }}
                onMouseDown={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                }}
                onClick={() => handleUnselect(value)}
              >
                <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
              </button>
            </Badge>
          ))}
          <input
            ref={inputRef}
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            className="ml-2 bg-transparent outline-none placeholder:text-muted-foreground flex-1"
            type="text"
          />
        </div>
      </div>
    </div>
  );
}
