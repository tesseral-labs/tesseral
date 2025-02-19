import React, { HTMLAttributes, PropsWithChildren } from 'react'
import { cva, VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const textDividerVariants = cva('block relative w-full', {
  variants: {
    variant: {
      default: 'my-4',
      tight: 'my-2',
      tighter: 'my-1',
      wide: 'my-6',
      wider: 'my-8',
      widest: 'my-12',
    },
    affects: {
      default:
        '[&>div.absolute>span]:border-muted-foreground [&>div.relative>span]:text-muted-foreground',
      muted:
        '[&>div.absolute>span]:border-muted [&>div.relative>span]:text-muted',
    },
  },
  defaultVariants: {
    variant: 'default',
    affects: 'default',
  },
})

export interface TextDividerProps
  extends HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof textDividerVariants>,
    PropsWithChildren {}

const TextDivider = React.forwardRef<HTMLDivElement, TextDividerProps>(
  ({ affects, children, className, variant, ...props }, ref) => {
    return (
      <div
        className={cn(textDividerVariants({ affects, variant, className }))}
        {...props}
        ref={ref}
      >
        <div className="absolute inset-0 flex items-center">
          <span className="w-full border-t"></span>
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-card px-2">{children}</span>
        </div>
      </div>
    )
  },
)

TextDivider.displayName = 'TextDivider'

export default TextDivider
