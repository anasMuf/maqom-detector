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

  // Parse Metadata if available
  let metadata: { title?: string; duration?: number; channel?: string } | null = null;
  try {
    if (data.metadata) {
      metadata = typeof data.metadata === 'string' ? JSON.parse(data.metadata) : data.metadata;
    }
  } catch (e) {
    console.error("Failed to parse metadata", e);
  }

  // Extract YouTube ID
  const getYoutubeId = (url: string) => {
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
    const match = url.match(regExp);
    return (match && match[2].length === 11) ? match[2] : null;
  };

  const youtubeId = data.input_type === 'youtube' ? getYoutubeId(data.input_source || '') : null;

  const getAudioUrl = () => {
    if (!data.input_source) return null;
    if (data.input_source.startsWith('http')) return data.input_source;
    return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/uploads/${data.input_source}`;
  };

  const audioUrl = (data.input_type === 'upload' || data.input_type === 'microphone') ? getAudioUrl() : null;

  return (
    <div className="flex flex-col gap-6 rounded-3xl border border-border bg-card p-6 shadow-sm overflow-hidden">
      {/* Media Player Section */}
      {(youtubeId || audioUrl) && (
        <div className="mb-2">
          {youtubeId ? (
            <div className="flex flex-col gap-3">
              <div className="aspect-video w-full rounded-2xl overflow-hidden bg-black border border-border shadow-md">
                <iframe
                  width="100%"
                  height="100%"
                  src={`https://www.youtube.com/embed/${youtubeId}`}
                  title="YouTube video player"
                  frameBorder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowFullScreen
                ></iframe>
              </div>
              {metadata?.title && (
                <div className="px-1">
                  <h3 className="text-sm font-bold text-foreground line-clamp-1">{metadata.title}</h3>
                  {metadata.channel && <p className="text-xs text-muted-foreground">{metadata.channel}</p>}
                </div>
              )}
            </div>
          ) : (
            <div className="p-4 bg-muted/50 rounded-2xl border border-border flex flex-col gap-3">
              <div className="flex items-center gap-3">
                <div className="h-10 w-10 rounded-full bg-brand-primary/10 flex items-center justify-center text-brand-primary">
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18V5l12-2v13" /><circle cx="6" cy="18" r="3" /><circle cx="18" cy="16" r="3" /></svg>
                </div>
                <div>
                  <h3 className="text-sm font-bold text-foreground">Sumber Audio</h3>
                  <p className="text-xs text-muted-foreground">Analisis dari {data.input_type === 'microphone' ? 'Rekaman Mikrofon' : 'File Upload'}</p>
                </div>
              </div>
              <audio controls className="w-full h-10">
                <source src={audioUrl!} type="audio/wav" />
                Browser Anda tidak mendukung elemen audio.
              </audio>
            </div>
          )}
        </div>
      )}

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
        <div className="mt-2">
          <h4 className="text-sm font-bold text-foreground mb-4">Kandidat Lainnya</h4>
          <div className="flex flex-col gap-5">
            {others.map((c) => (
              <div key={c.maqam_id} className="flex flex-col gap-2">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-foreground">{c.name_latin}</span>
                  <span className="text-brand-primary font-bold">{Math.round((c.confidence_score || 0) * 100)}%</span>
                </div>
                <ProgressBar progress={(c.confidence_score || 0) * 100} showLabel={false} className="h-2" />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
