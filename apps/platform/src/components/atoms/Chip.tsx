import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '../../lib/utils';

const chipVariants = cva(
  'inline-flex items-center justify-center rounded-full text-sm font-medium transition-colors cursor-default',
  {
    variants: {
      variant: {
        filled: 'bg-brand-primary-subtle text-brand-primary',
        outline: 'border border-brand-primary text-brand-primary',
      },
      size: {
        sm: 'px-2 py-1 text-xs',
        md: 'px-3 py-1.5',
      },
    },
    defaultVariants: {
      variant: 'filled',
      size: 'md',
    },
  }
);

export interface ChipProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof chipVariants> {}

function Chip({ className, variant, size, ...props }: ChipProps) {
  return (
    <div
      className={cn(chipVariants({ variant, size }), className)}
      {...props}
    />
  );
}

export { Chip, chipVariants };
