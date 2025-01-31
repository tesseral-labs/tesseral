import * as React from 'react'
import { cn } from '@/lib/utils'

export const DetailsGrid = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn('text-sm grid grid-flow-col auto-cols-fr gap-x-2', className)}
    {...props}
  />
))
DetailsGrid.displayName = 'DetailsGrid'

export const DetailsGridColumn = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      'border-r last:border-r-0 px-4 first:pl-0 last:pr-0 border-gray-200 flex flex-col gap-y-3',
      className,
    )}
    {...props}
  />
))
DetailsGridColumn.displayName = 'DetailsGridColumn'

export const DetailsGridEntry = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn('', className)} {...props} />
))
DetailsGridEntry.displayName = 'DetailsGridEntry'

export const DetailsGridKey = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn('font-semibold', className)} {...props} />
))
DetailsGridKey.displayName = 'DetailsGridKey'

export const DetailsGridValue = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn('truncate', className)} {...props} />
))
DetailsGridValue.displayName = 'DetailsGridValue'
