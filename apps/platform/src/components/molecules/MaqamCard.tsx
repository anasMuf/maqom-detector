import * as React from 'react';
import { type DtoMaqamListResponse } from '../../api/model';
import { Badge } from '../atoms/Badge';
import { cn } from '../../lib/utils';

interface MaqamCardProps extends React.HTMLAttributes<HTMLDivElement> {
  maqam: DtoMaqamListResponse;
}

export function MaqamCard({ maqam, className, ...props }: MaqamCardProps) {
  return (
    <div
      className={cn(
        'flex flex-col gap-3 rounded-2xl border border-border bg-card p-5 shadow-sm transition-shadow hover:shadow-md cursor-pointer',
        className
      )}
      {...props}
    >
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-xl font-bold text-card-foreground">
            {maqam.name_latin}
          </h3>
          <p className="text-sm text-muted-foreground mt-1">
            Interval: <span className="font-mono">{maqam.interval_description}</span>
          </p>
        </div>
        <span className="font-arabic text-3xl font-bold text-brand-primary">
          {maqam.name_arabic}
        </span>
      </div>

      {maqam.emotion_tags && maqam.emotion_tags.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-2">
          {maqam.emotion_tags.map((tag) => (
            <Badge key={tag} variant="neutral" className="bg-gray-100 text-gray-700 hover:bg-gray-200">
              {tag}
            </Badge>
          ))}
        </div>
      )}
    </div>
  );
}
