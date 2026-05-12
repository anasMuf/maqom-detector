import { createFileRoute, Link } from '@tanstack/react-router';
import { useGetHistory, useDeleteHistoryId } from '../api/endpoints/history/history';
import { Badge } from '../components/atoms/Badge';
import { buttonVariants } from '../components/atoms/Button';
import { Spinner } from '../components/atoms/Spinner';
import { Trash2, Music, History as HistoryIcon } from 'lucide-react';

export const Route = createFileRoute('/history')({
  component: HistoryPage,
});

function HistoryPage() {
  const { data, isLoading, refetch } = useGetHistory({ page: 1, limit: 50 });
  const deleteMutation = useDeleteHistoryId();

  const handleDelete = (id?: string) => {
    if (!id) return;
    if (confirm('Hapus riwayat analisis ini?')) {
      deleteMutation.mutate({ id }, {
        onSuccess: () => refetch()
      });
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    const d = new Date(dateString);
    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(d);
  };

  return (
    <div className="mx-auto max-w-4xl flex flex-col gap-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Riwayat Analisis</h1>
          <p className="text-muted-foreground mt-2">Daftar deteksi maqam yang pernah Anda lakukan di sesi ini.</p>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-20">
          <Spinner size="lg" />
        </div>
      ) : data?.data?.items && data.data.items.length > 0 ? (
        <div className="flex flex-col gap-4">
          {data.data.items.map((item) => (
            <div key={item.id} className="flex flex-col sm:flex-row sm:items-center justify-between p-4 rounded-xl border border-border bg-card shadow-sm gap-4 transition-colors hover:bg-muted/20">
              <div className="flex items-start gap-4">
                <div className="mt-1 flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-brand-primary-subtle text-brand-primary">
                  <Music size={18} />
                </div>
                <div>
                  <h3 className="font-bold text-foreground flex items-center gap-2">
                    {item.status === 'completed' ? item.maqam_name_latin : 'Menunggu Analisis'}
                    {item.status === 'completed' && (
                      <span className="text-sm font-normal text-muted-foreground">
                        ({Math.round((item.confidence_score || 0) * 100)}%)
                      </span>
                    )}
                  </h3>
                  <div className="flex flex-wrap items-center gap-2 mt-1.5">
                    <Badge variant={item.status === 'completed' ? 'success' : item.status === 'failed' ? 'error' : 'warning'} className="text-[10px] px-1.5 py-0">
                      {item.status}
                    </Badge>
                    <span className="text-xs text-muted-foreground uppercase">{item.input_type}</span>
                    <span className="text-xs text-muted-foreground">•</span>
                    <span className="text-xs text-muted-foreground">{formatDate(item.created_at)}</span>
                  </div>
                </div>
              </div>
              
              <div className="flex items-center gap-2 sm:self-center self-end">
                {item.status === 'completed' && (
                  <Link to="/analyze" search={{ id: item.id }} className={buttonVariants({ variant: 'outline', size: 'sm' })}>Detail</Link>
                )}
                <button 
                  className={buttonVariants({ variant: 'ghost', size: 'icon' }) + " text-muted-foreground hover:text-red-500 hover:bg-red-50"}
                  onClick={() => handleDelete(item.id)}
                  disabled={deleteMutation.isPending}
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center p-12 text-center border border-dashed border-border rounded-2xl bg-card">
          <HistoryIcon size={48} className="text-muted-foreground/50 mb-4" />
          <h3 className="text-lg font-bold text-foreground">Belum ada riwayat</h3>
          <p className="text-muted-foreground mt-2 max-w-sm">Anda belum melakukan analisis maqam apapun pada sesi ini.</p>
          <Link to="/analyze" className={buttonVariants({ variant: 'filled' }) + " mt-6"}>Mulai Deteksi Sekarang</Link>
        </div>
      )}
    </div>
  );
}
