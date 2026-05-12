import { createFileRoute } from '@tanstack/react-router';
import { useGetMaqamat } from '../api/endpoints/maqam/maqam';
import { MaqamCard } from '../components/molecules/MaqamCard';
import { Spinner } from '../components/atoms/Spinner';

export const Route = createFileRoute('/maqamat')({
  component: MaqamatPage,
});

function MaqamatPage() {
  const { data, isLoading, error } = useGetMaqamat();

  return (
    <div className="mx-auto max-w-5xl flex flex-col gap-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-foreground">Kamus Maqamat</h1>
        <p className="text-muted-foreground mt-2 max-w-2xl mx-auto">
          Pelajari 8 maqam dasar dalam musik Arab yang sering digunakan dalam kesenian banjari dan tilawah.
        </p>
      </div>

      {isLoading && (
        <div className="flex justify-center py-20">
          <Spinner size="lg" />
        </div>
      )}

      {!!error && (
        <div className="p-4 bg-red-50 text-red-600 rounded-xl text-center">
          Gagal memuat data maqam. Silakan coba lagi.
        </div>
      )}

      {data?.data && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {data.data.map((maqam) => (
            <MaqamCard key={maqam.id} maqam={maqam} />
          ))}
        </div>
      )}
    </div>
  );
}
