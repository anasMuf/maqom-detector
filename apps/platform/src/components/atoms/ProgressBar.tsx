import * as React from 'react';

interface ProgressBarProps {
  progress: number; // 0 to 100
  label?: string;
  showLabel?: boolean;
  className?: string;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({ 
  progress, 
  label, 
  showLabel = true, 
  className = '' 
}) => {
  // Ensure progress is between 0 and 100
  const clampedProgress = Math.min(Math.max(progress, 0), 100);

  return (
    <div className={`w-full flex flex-col gap-1.5 ${className}`}>
      {showLabel && (
        <div className="flex justify-between items-center text-sm font-medium">
          <span className="text-muted-foreground">{label || 'Memproses...'}</span>
          <span className="text-brand-primary">{Math.round(clampedProgress)}%</span>
        </div>
      )}
      <div className="h-3 w-full bg-muted rounded-full overflow-hidden border border-border shadow-inner">
        <div
          className="h-full bg-brand-primary transition-all duration-500 ease-out rounded-full shadow-[0_0_10px_rgba(var(--brand-primary-rgb),0.5)]"
          style={{ width: `${clampedProgress}%` }}
        >
          {/* Shimmer effect */}
          <div className="w-full h-full relative overflow-hidden">
            <div className="absolute inset-0 bg-linier-to-r from-transparent via-white/20 to-transparent translate-x-full animate-[shimmer_2s_infinite]"></div>
          </div>
        </div>
      </div>
    </div>
  );
};
