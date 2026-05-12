import * as React from 'react';
import { cn } from '../../lib/utils';

export interface SkeletonLoaderProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'text' | 'image' | 'card' | 'circular';
}

export function SkeletonLoader({
  className,
  variant = 'text',
  ...props
}: SkeletonLoaderProps) {
  return (
    <div
      className={cn(
        'animate-pulse rounded-md bg-muted',
        {
          'h-4 w-full': variant === 'text',
          'h-full w-full': variant === 'image',
          'h-32 w-full': variant === 'card',
          'h-12 w-12 rounded-full': variant === 'circular',
        },
        className
      )}
      {...props}
    />
  );
}
