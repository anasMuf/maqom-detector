import { createFileRoute, Link } from '@tanstack/react-router';
import { Music, Upload, Mic, History } from 'lucide-react';
import { buttonVariants } from '../components/atoms/Button';

export const Route = createFileRoute('/')({
  component: Index,
});

function Index() {
  return (
    <div className="flex flex-col items-center justify-center py-10 md:py-20 text-center">
      <div className="mb-6 rounded-full bg-brand-primary-subtle p-4 text-brand-primary">
        <Music size={48} />
      </div>
      
      <h1 className="mb-4 text-4xl font-extrabold tracking-tight md:text-6xl text-foreground">
        Deteksi <span className="text-brand-primary">Maqam</span> Otomatis
      </h1>
      
      <p className="mb-8 max-w-2xl text-lg text-muted-foreground md:text-xl">
        Unggah audio, rekam suara, atau masukkan link YouTube untuk mendeteksi maqam secara real-time. Didesain khusus untuk musisi banjari dan pecinta musik Arab.
      </p>

      <div className="flex flex-col gap-4 sm:flex-row mb-16">
        <Link to="/analyze" className={buttonVariants({ size: 'lg' }) + " gap-2"}>
          <Mic size={20} />
          Mulai Deteksi
        </Link>
        <Link to="/maqamat" className={buttonVariants({ variant: 'outline', size: 'lg' }) + " gap-2"}>
          <Music size={20} />
          Pelajari Maqamat
        </Link>
      </div>

      {/* Feature Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 w-full max-w-4xl text-left">
        <div className="flex flex-col gap-3 rounded-2xl bg-card p-6 border border-border shadow-sm">
          <div className="w-12 h-12 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center">
            <Upload size={24} />
          </div>
          <h3 className="text-xl font-bold">Upload Audio</h3>
          <p className="text-muted-foreground text-sm">Upload rekaman dari grup banjari Anda (MP3/WAV) dan deteksi polanya.</p>
        </div>
        <div className="flex flex-col gap-3 rounded-2xl bg-card p-6 border border-border shadow-sm">
          <div className="w-12 h-12 rounded-lg bg-green-100 text-green-600 flex items-center justify-center">
            <Mic size={24} />
          </div>
          <h3 className="text-xl font-bold">Rekam Langsung</h3>
          <p className="text-muted-foreground text-sm">Gunakan mikrofon untuk bersenandung atau menyanyi, AI akan mencocokkan maqamnya.</p>
        </div>
        <div className="flex flex-col gap-3 rounded-2xl bg-card p-6 border border-border shadow-sm">
          <div className="w-12 h-12 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center">
            <History size={24} />
          </div>
          <h3 className="text-xl font-bold">Riwayat Deteksi</h3>
          <p className="text-muted-foreground text-sm">Simpan dan akses kembali hasil analisis sebelumnya tanpa harus login.</p>
        </div>
      </div>
    </div>
  );
}
