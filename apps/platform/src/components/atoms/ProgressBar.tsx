import * as React from 'react';
import { cn } from '../../lib/utils';

interface ProgressBarProps extends React.HTMLAttributes<HTMLDivElement> {
  value: number; // 0 to 100
  variant?: 'brand' | 'success' | 'warning' | 'error';
}

const variantStyles = {
  brand: 'bg-brand-primary',
  success: 'bg-confidence-high',
  warning: 'bg-confidence-medium',
  error: 'bg-confidence-very-low',
};

export function ProgressBar({
  value,
  variant = 'brand',
  className,
  ...props
}: ProgressBarProps) {
  const safeValue = Math.min(100, Math.max(0, value));

  return (
    <div
      className={cn('h-2 w-full overflow-hidden rounded-full bg-secondary', className)}
      {...props}
      role="progressbar"
      aria-valuenow={safeValue}
      aria-valuemin={0}
      aria-valuemax={100}
    >
      <div
        className={cn('h-full w-full flex-1 transition-all duration-500 ease-in-out', variantStyles[variant])}
        style={{ transform: `translateX(-${100 - safeValue}%)` }}
      />
    </div>
  );
}
