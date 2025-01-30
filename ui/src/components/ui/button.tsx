import React, { CSSProperties, useEffect, useState } from 'react'
import { cva, type VariantProps } from 'class-variance-authority'

import { cn, isColorDark } from '@/lib/utils'
import useSettings from '@/lib/settings'
import useDarkMode from '@/lib/dark-mode'

const buttonVariants = cva(
  'inline-flex font-semibold items-center justify-center whitespace-nowrap rounded-md text-sm ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none',
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground hover:bg-primary/90 disabled:text-muted-foreground disabled:bg-muted',
        destructive:
          'bg-destructive text-destructive-foreground hover:bg-destructive/90',
        outline:
          'border border-border bg-background hover:bg-accent hover:text-accent-foreground',
        secondary:
          'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        ghost: 'hover:bg-accent hover:text-accent-foreground',
        link: 'text-primary underline-offset-4 hover:underline',
      },
      size: {
        default: 'h-10 px-4 py-2',
        sm: 'h-9 rounded-md px-3',
        lg: 'h-11 rounded-md px-8',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const isDarkMode = useDarkMode()
    const [customStyle, setCustomStyle] = useState<CSSProperties>(
      props.style ? { ...props.style } : {},
    )

    const settings = useSettings()

    useEffect(() => {
      const style: CSSProperties = { ...customStyle }

      if (
        settings?.darkModePrimaryColor &&
        isDarkMode &&
        settings?.detectDarkModeEnabled &&
        (!variant || variant === 'default')
      ) {
        style.backgroundColor = settings?.darkModePrimaryColor
        style.color = isColorDark(settings?.primaryColor) ? 'white' : 'black'
      } else if (
        settings?.primaryColor &&
        (!variant || variant === 'default')
      ) {
        style.backgroundColor = settings?.primaryColor
        style.color = isColorDark(settings?.primaryColor) ? 'white' : 'black'
      }

      setCustomStyle(style)
    }, [settings, variant])

    return (
      <button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        style={customStyle}
        {...props}
      />
    )
  },
)

export { Button, buttonVariants }
