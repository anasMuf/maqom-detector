import * as React from 'react';
import { type DtoAnalysisDetailResponse } from '../../api/model';
import { Badge } from '../atoms/Badge';
import { ProgressBar } from '../atoms/ProgressBar';

interface AnalysisResultProps {
  data: DtoAnalysisDetailResponse;
}

export function AnalysisResult({ data }: AnalysisResultProps) {
  if (!data.candidates || data.candidates.length === 0) return null;

  const topMaqam = data.candidates[0];
  const others = data.candidates.slice(1);

  return (
    <div className="flex flex-col gap-6 rounded-3xl border border-border bg-card p-6 shadow-sm">
      {/* Top Result */}
      <div className="flex flex-col items-center text-center">
        <span className="text-sm font-medium text-brand-primary mb-1 uppercase tracking-wider">
          Maqam Terdeteksi
        </span>
        <h2 className="text-4xl font-bold text-foreground">
          {topMaqam.name_latin}
        </h2>
        <span className="font-arabic text-4xl mt-2 text-brand-secondary">
          {topMaqam.name_arabic}
        </span>
        
        <div className="mt-4 flex items-center gap-2">
          <Badge variant={data.confidence_label === 'tinggi' || data.confidence_label === 'sangat_tinggi' ? 'success' : 'warning'}>
            Keyakinan {data.confidence_label?.replace('_', ' ')}
          </Badge>
          <span className="text-sm text-muted-foreground font-mono">
            {Math.round((topMaqam.confidence_score || 0) * 100)}%
          </span>
        </div>
      </div>

      <hr className="border-border" />

      {/* Explanation */}
      {data.explanation_text && (
        <div>
          <h4 className="text-sm font-bold text-foreground mb-2">Penjelasan</h4>
          <div className="text-sm text-muted-foreground leading-relaxed whitespace-pre-wrap">
            {data.explanation_text}
          </div>
        </div>
      )}

      {/* Other Candidates */}
      {others.length > 0 && (
        <div>
          <h4 className="text-sm font-bold text-foreground mb-3">Kandidat Lainnya</h4>
          <div className="flex flex-col gap-3">
            {others.map((c) => (
              <div key={c.maqam_id} className="flex flex-col gap-1.5">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-foreground">{c.name_latin}</span>
                  <span className="text-muted-foreground">{Math.round((c.confidence_score || 0) * 100)}%</span>
                </div>
                <ProgressBar value={(c.confidence_score || 0) * 100} variant="brand" className="h-1.5" />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
