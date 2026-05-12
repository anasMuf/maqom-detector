import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../../lib/utils';

const spinnerVariants = cva('animate-spin inline-block rounded-full', {
  variants: {
    size: {
      sm: 'w-4 h-4 border-2',
      md: 'w-6 h-6 border-2',
      lg: 'w-8 h-8 border-3',
    },
    color: {
      brand: 'border-brand-primary border-t-transparent',
      white: 'border-white border-t-transparent',
      neutral: 'border-gray-500 border-t-transparent',
    },
  },
  defaultVariants: {
    size: 'md',
    color: 'brand',
  },
});

export interface SpinnerProps
  extends Omit<React.HTMLAttributes<HTMLDivElement>, 'color'>,
    VariantProps<typeof spinnerVariants> {}

export function Spinner({ className, size, color, ...props }: SpinnerProps) {
  return (
    <div
      className={cn(spinnerVariants({ size, color, className }))}
      role="status"
      aria-label="loading"
      {...props}
    >
      <span className="sr-only">Loading...</span>
    </div>
  );
}
