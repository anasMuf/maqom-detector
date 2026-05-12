import * as React from 'react';
import { Upload, FileAudio, X } from 'lucide-react';
import { Button } from '../atoms/Button';
import { cn } from '../../lib/utils';

export interface AudioUploaderProps {
  onFileSelect: (file: File | null) => void;
  selectedFile: File | null;
  isUploading?: boolean;
  className?: string;
}

export function AudioUploader({
  onFileSelect,
  selectedFile,
  isUploading,
  className,
}: AudioUploaderProps) {
  const inputRef = React.useRef<HTMLInputElement>(null);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      onFileSelect(e.dataTransfer.files[0]);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      onFileSelect(e.target.files[0]);
    }
  };

  return (
    <div className={cn('w-full', className)}>
      {!selectedFile ? (
        <div
          onDragOver={handleDragOver}
          onDrop={handleDrop}
          onClick={() => inputRef.current?.click()}
          className="flex cursor-pointer flex-col items-center justify-center rounded-2xl border-2 border-dashed border-border bg-muted/30 py-12 px-4 transition-colors hover:bg-muted/50"
        >
          <div className="rounded-full bg-brand-primary-subtle p-3 text-brand-primary mb-4">
            <Upload size={24} />
          </div>
          <p className="text-center font-medium text-foreground">
            Klik atau drop file audio di sini
          </p>
          <p className="mt-1 text-center text-sm text-muted-foreground">
            MP3, WAV, M4A up to 50MB
          </p>
          <input
            ref={inputRef}
            type="file"
            accept="audio/*"
            className="hidden"
            onChange={handleChange}
          />
        </div>
      ) : (
        <div className="flex items-center justify-between rounded-xl border border-border bg-card p-4">
          <div className="flex items-center gap-3 overflow-hidden">
            <div className="rounded-full bg-brand-primary-subtle p-2 text-brand-primary shrink-0">
              <FileAudio size={20} />
            </div>
            <div className="overflow-hidden">
              <p className="truncate font-medium text-sm text-foreground">
                {selectedFile.name}
              </p>
              <p className="text-xs text-muted-foreground">
                {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
              </p>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onFileSelect(null)}
            disabled={isUploading}
            className="shrink-0 text-muted-foreground hover:text-destructive"
          >
            <X size={18} />
          </Button>
        </div>
      )}
    </div>
  );
}
