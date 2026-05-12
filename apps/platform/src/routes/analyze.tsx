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
import { Spinner } from '../components/atoms/Spinner';

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

      <div className="flex justify-center gap-2 mb-2 p-1 bg-muted rounded-xl">
        <button
          onClick={() => setActiveTab('upload')}
          className={`flex-1 rounded-lg py-2 text-sm font-medium transition-colors ${activeTab === 'upload' ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'}`}
        >
          Upload
        </button>
        <button
          onClick={() => setActiveTab('youtube')}
          className={`flex-1 rounded-lg py-2 text-sm font-medium transition-colors ${activeTab === 'youtube' ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'}`}
        >
          YouTube
        </button>
        <button
          onClick={() => setActiveTab('record')}
          className={`flex-1 rounded-lg py-2 text-sm font-medium transition-colors ${activeTab === 'record' ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'}`}
        >
          Mikrofon
        </button>
      </div>

      <div className="bg-card border border-border rounded-2xl p-6 shadow-sm">
        {activeTab === 'upload' && (
          <div className="flex flex-col gap-6">
            <AudioUploader onFileSelect={setFile} selectedFile={file} isUploading={isProcessing} />
            <Button onClick={handleUploadSubmit} disabled={!file || isProcessing} isLoading={uploadMutation.isPending} className="w-full">
              Analisis Audio
            </Button>
          </div>
        )}

        {activeTab === 'youtube' && (
          <div className="flex flex-col gap-6">
            <div>
              <label className="block text-sm font-medium text-foreground mb-2">Link YouTube</label>
              <Input
                placeholder="https://www.youtube.com/watch?v=..."
                value={youtubeUrl}
                onChange={(e) => setYoutubeUrl(e.target.value)}
                disabled={isProcessing}
              />
              <p className="text-xs text-muted-foreground mt-2">Sistem akan mengambil 120 detik pertama dari video.</p>
            </div>
            <Button onClick={handleYoutubeSubmit} disabled={!youtubeUrl || isProcessing} isLoading={youtubeMutation.isPending} className="w-full">
              Analisis YouTube
            </Button>
          </div>
        )}

        {activeTab === 'record' && (
          <div className="flex flex-col gap-6 items-center py-6">
            {isRecording ? (
              <div className="flex flex-col items-center gap-4">
                <div className="h-24 w-24 rounded-full bg-red-100 flex items-center justify-center animate-pulse">
                  <div className="h-12 w-12 rounded-full bg-red-500"></div>
                </div>
                <p className="text-sm font-medium text-red-500">Sedang merekam...</p>
                <Button onClick={stopRecording} variant="outline" className="text-red-500 border-red-500 hover:bg-red-50">
                  Berhenti & Analisis
                </Button>
              </div>
            ) : (
              <div className="flex flex-col items-center gap-4">
                <div className="h-24 w-24 rounded-full bg-brand-primary-subtle flex items-center justify-center text-brand-primary">
                  <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"/><path d="M19 10v2a7 7 0 0 1-14 0v-2"/><line x1="12" x2="12" y1="19" y2="22"/></svg>
                </div>
                <p className="text-sm font-medium text-muted-foreground">Izinkan akses mikrofon untuk mulai</p>
                <Button onClick={startRecording} disabled={isProcessing} isLoading={recordMutation.isPending}>
                  Mulai Rekaman
                </Button>
              </div>
            )}
          </div>
        )}
      </div>

      {isProcessing && analysisStatus !== 'completed' && analysisStatus !== 'failed' && (
        <div className="flex flex-col items-center justify-center p-12 border border-border rounded-2xl bg-card">
          <Spinner size="lg" className="mb-4" />
          <p className="text-lg font-medium text-foreground">Menganalisis Audio...</p>
          <p className="text-sm text-muted-foreground mt-2">Ini mungkin memakan waktu hingga 30 detik</p>
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
