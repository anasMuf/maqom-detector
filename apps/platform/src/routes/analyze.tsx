import * as React from 'react';
import { createFileRoute } from '@tanstack/react-router';
import {
  usePostAnalyzeUpload,
  usePostAnalyzeYoutube,
  usePostAnalyzeRecord,
  useGetAnalysesId
} from '../api/endpoints/analyze/analyze';
import { AudioUploader } from '../components/molecules/AudioUploader';
import { AnalysisResult } from '../components/molecules/AnalysisResult';
import { Button } from '../components/atoms/Button';
import { Input } from '../components/atoms/Input';
import { ProgressBar } from '../components/atoms/ProgressBar';

export const Route = createFileRoute('/analyze')({
  component: AnalyzePage,
  validateSearch: (search: Record<string, unknown>): { id?: string } => {
    return {
      id: search.id as string | undefined,
    }
  },
});

function AnalyzePage() {
  const { id: searchId } = Route.useSearch();
  const [activeTab, setActiveTab] = React.useState<'upload' | 'youtube' | 'record'>('upload');

  // Simulated progress logic
  const [progress, setProgress] = React.useState(0);

  // const { id: idFromState } = Route.useSearch(); // To keep track of state changes

  const [file, setFile] = React.useState<File | null>(null);
  const [youtubeUrl, setYoutubeUrl] = React.useState('');
  const [isRecording, setIsRecording] = React.useState(false);
  const [mediaRecorder, setMediaRecorder] = React.useState<MediaRecorder | null>(null);

  const [pollingId, setPollingId] = React.useState<string | null>(searchId || null);

  // Mutations
  const uploadMutation = usePostAnalyzeUpload();
  const youtubeMutation = usePostAnalyzeYoutube();
  const recordMutation = usePostAnalyzeRecord();

  // Polling query
  const { data: analysisRaw } = useGetAnalysesId(
    pollingId as string,
    {
      query: {
        enabled: !!pollingId,
        refetchInterval: (query) => {
          const raw = query.state.data as Record<string, unknown> | undefined;
          if (!raw?.data) return 2000;
          const inner = raw.data as Record<string, unknown>;
          const st = inner?.status as string | undefined;
          if (st === 'completed' || st === 'failed') return false;
          return 2000;
        },
      }
    }
  );

  // Extract the inner data from API envelope { data: {...}, success: true }
  const analysisEnvelope = analysisRaw as unknown as { data?: Record<string, unknown>; success?: boolean } | undefined;
  const analysisDetail = analysisEnvelope?.data as import('../api/model').DtoAnalysisDetailResponse | undefined;
  const analysisStatus = analysisDetail?.status ?? null;

  const isProcessing = uploadMutation.isPending || youtubeMutation.isPending || recordMutation.isPending || analysisStatus === 'pending' || analysisStatus === 'processing';

  React.useEffect(() => {
    let interval: NodeJS.Timeout;

    if (isProcessing) {
      interval = setInterval(() => {
        setProgress((prev) => {
          if (prev >= 95) return prev;
          const increment = prev < 30 ? 3 : prev < 70 ? 1 : 0.3;
          return prev + increment;
        });
      }, 500);
    } else if (analysisStatus === 'completed') {
      setProgress(100);
    } else {
      setProgress(0);
    }

    return () => clearInterval(interval);
  }, [isProcessing, analysisStatus]);

  // Helper: extract analysis_id from mutation response envelope
  const extractAnalysisId = (res: unknown): string | undefined => {
    const envelope = res as { data?: { analysis_id?: string } };
    return envelope?.data?.analysis_id;
  };

  const handleUploadSubmit = () => {
    if (!file) return;
    uploadMutation.mutate({ data: { file } }, {
      onSuccess: (res) => {
        const id = extractAnalysisId(res);
        if (id) setPollingId(id);
      }
    });
  };

  const handleYoutubeSubmit = () => {
    if (!youtubeUrl) return;
    youtubeMutation.mutate({ data: { url: youtubeUrl, segment_start: 0, segment_duration: 120 } }, {
      onSuccess: (res) => {
        const id = extractAnalysisId(res);
        if (id) setPollingId(id);
      }
    });
  };

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new MediaRecorder(stream);
      const chunks: BlobPart[] = [];

      recorder.ondataavailable = (e) => chunks.push(e.data);
      recorder.onstop = () => {
        const blob = new Blob(chunks, { type: 'audio/wav' });
        const recordFile = new File([blob], 'recording.wav', { type: 'audio/wav' });
        recordMutation.mutate({ data: { file: recordFile, mode: 'microphone' } }, {
          onSuccess: (res) => {
            const id = extractAnalysisId(res);
            if (id) setPollingId(id);
          }
        });
      };

      recorder.start();
      setMediaRecorder(recorder);
      setIsRecording(true);
    } catch (err) {
      console.error('Failed to access microphone', err);
      alert('Gagal mengakses mikrofon');
    }
  };

  const stopRecording = () => {
    if (mediaRecorder) {
      mediaRecorder.stop();
      mediaRecorder.stream.getTracks().forEach(track => track.stop());
      setIsRecording(false);
    }
  };

  return (
    <div className="mx-auto max-w-3xl flex flex-col gap-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-foreground">Deteksi Maqam</h1>
        <p className="text-muted-foreground mt-2">Pilih sumber audio yang ingin dianalisis</p>
      </div>

      <div className="flex justify-center p-1.5 bg-muted/50 border border-border/50 rounded-2xl mb-2 backdrop-blur-sm">
        <button
          onClick={() => setActiveTab('upload')}
          className={`flex-1 flex items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold transition-all duration-200 ${activeTab === 'upload'
              ? 'bg-background text-brand-primary shadow-sm ring-1 ring-black/5'
              : 'text-muted-foreground hover:text-foreground hover:bg-black/5'
            }`}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" /><polyline points="17 8 12 3 7 8" /><line x1="12" x2="12" y1="3" y2="15" /></svg>
          Upload
        </button>
        <button
          onClick={() => setActiveTab('youtube')}
          className={`flex-1 flex items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold transition-all duration-200 ${activeTab === 'youtube'
              ? 'bg-background text-brand-primary shadow-sm ring-1 ring-black/5'
              : 'text-muted-foreground hover:text-foreground hover:bg-black/5'
            }`}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M22.54 6.42a2.78 2.78 0 0 0-1.94-2C18.88 4 12 4 12 4s-6.88 0-8.6.46a2.78 2.78 0 0 0-1.94 2A29 29 0 0 0 1 11.75a29 29 0 0 0 .46 5.33A2.78 2.78 0 0 0 3.4 19c1.72.46 8.6.46 8.6.46s6.88 0 8.6-.46a2.78 2.78 0 0 0 1.94-2 29 29 0 0 0 .46-5.25 29 29 0 0 0-.46-5.33z" /><polygon points="9.75 15.02 15.5 11.75 9.75 8.48 9.75 15.02" /></svg>
          YouTube
        </button>
        <button
          onClick={() => setActiveTab('record')}
          className={`flex-1 flex items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold transition-all duration-200 ${activeTab === 'record'
              ? 'bg-background text-brand-primary shadow-sm ring-1 ring-black/5'
              : 'text-muted-foreground hover:text-foreground hover:bg-black/5'
            }`}
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3z" /><path d="M19 10v2a7 7 0 0 1-14 0v-2" /><line x1="12" x2="12" y1="19" y2="22" /></svg>
          Mikrofon
        </button>
      </div>

      <div className="bg-card border border-border/60 rounded-3xl p-8 shadow-xl shadow-black/3 backdrop-blur-sm">
        {activeTab === 'upload' && (
          <div className="flex flex-col gap-8">
            <AudioUploader onFileSelect={setFile} selectedFile={file} isUploading={isProcessing} />
            <Button onClick={handleUploadSubmit} disabled={!file || isProcessing} isLoading={uploadMutation.isPending} className="w-full h-12 text-base font-bold rounded-2xl shadow-lg shadow-brand-primary/20">
              Analisis Audio
            </Button>
          </div>
        )}

        {activeTab === 'youtube' && (
          <div className="flex flex-col gap-8">
            <div>
              <label className="block text-sm font-bold text-foreground mb-3">Link YouTube</label>
              <Input
                placeholder="https://www.youtube.com/watch?v=..."
                value={youtubeUrl}
                onChange={(e) => setYoutubeUrl(e.target.value)}
                disabled={isProcessing}
                className="h-12 rounded-xl"
              />
              <p className="text-xs text-muted-foreground mt-3 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="10" /><line x1="12" x2="12" y1="16" y2="12" /><line x1="12" x2="12.01" y1="8" y2="8" /></svg>
                Sistem akan mengambil 120 detik pertama dari video.
              </p>
            </div>
            <Button onClick={handleYoutubeSubmit} disabled={!youtubeUrl || isProcessing} isLoading={youtubeMutation.isPending} className="w-full h-12 text-base font-bold rounded-2xl shadow-lg shadow-brand-primary/20">
              Analisis YouTube
            </Button>
          </div>
        )}

        {activeTab === 'record' && (
          <div className="flex flex-col items-center gap-8 py-4">
            <div className={`h-28 w-28 rounded-full flex items-center justify-center transition-all duration-500 ${isRecording ? 'bg-red-500/10 text-red-500 animate-pulse ring-4 ring-red-500/5' : 'bg-brand-primary-subtle text-brand-primary shadow-inner'}`}>
              <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
                <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z" />
                <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
                <line x1="12" x2="12" y1="19" y2="22" />
              </svg>
            </div>
            <div className="text-center">
              <h3 className="text-xl font-bold text-foreground">{isRecording ? 'Sedang Merekam...' : 'Siap Merekam'}</h3>
              <p className="text-sm text-muted-foreground mt-2">Pastikan suara terdengar jelas untuk akurasi terbaik</p>
            </div>
            <Button
              onClick={isRecording ? stopRecording : startRecording}
              disabled={isProcessing}
              className={`w-full h-12 text-base font-bold rounded-2xl shadow-lg transition-all ${isRecording
                  ? 'bg-red-500 hover:bg-red-600 text-white shadow-red-500/20'
                  : 'bg-brand-primary hover:bg-brand-primary-hover text-white shadow-brand-primary/20'
                }`}
            >
              {isRecording ? 'Berhenti & Analisis' : 'Mulai Rekaman'}
            </Button>
          </div>
        )}
      </div>

      {isProcessing && analysisStatus !== 'completed' && analysisStatus !== 'failed' && (
        <div className="flex flex-col items-center justify-center p-12 border border-border rounded-2xl bg-card shadow-sm">
          <ProgressBar progress={progress} label="Menganalisis Maqam..." className="max-w-md w-full" />
          <p className="text-sm text-muted-foreground mt-6 text-center italic">
            {progress < 30 ? 'Menyiapkan audio...' : progress < 70 ? 'Mengekstrak pola nada (PCP)...' : 'Mencocokkan dengan referensi maqam...'}
          </p>
        </div>
      )}

      {analysisStatus === 'completed' && analysisDetail && (
        <div className="mt-4">
          <h2 className="text-xl font-bold mb-4">Hasil Analisis</h2>
          <AnalysisResult data={analysisDetail} />
        </div>
      )}

      {analysisStatus === 'failed' && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-xl text-red-600">
          <p className="font-bold">Analisis Gagal</p>
          <p className="text-sm mt-1">{analysisDetail?.error_message || 'Terjadi kesalahan saat memproses audio.'}</p>
        </div>
      )}
    </div>
  );
}
