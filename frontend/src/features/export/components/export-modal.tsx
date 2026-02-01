import { useState, useCallback, type ReactNode } from 'react';
import { Button } from '@/components/ui/button';
import type { Keyword } from '@/lib/schemas';
import api from '@/lib/axios';

interface FilterSummary {
  keyword_ids: string[];
  start_date: string;
  end_date: string;
  search: string;
  new_only: boolean;
}

interface ExportModalProps {
  isOpen: boolean;
  onClose: () => void;
  filters: FilterSummary;
  keywords: Keyword[];
}

type ModalState = 'idle' | 'exporting' | 'rate_limited' | 'error';

interface ExportError {
  code: string;
  message: string;
  retryAfter?: number;
}

function FilterChip({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700 border border-blue-200">
      {children}
    </span>
  );
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
}

export function ExportModal({ isOpen, onClose, filters, keywords }: ExportModalProps) {
  const [state, setState] = useState<ModalState>('idle');
  const [error, setError] = useState<ExportError | null>(null);

  const handleExport = useCallback(async () => {
    setState('exporting');
    setError(null);

    try {
      const params = new URLSearchParams();
      if (filters.keyword_ids.length > 0) {
        params.set('keyword_ids', filters.keyword_ids.join(','));
      }
      if (filters.start_date) params.set('start_date', filters.start_date);
      if (filters.end_date) params.set('end_date', filters.end_date);
      if (filters.search) params.set('search', filters.search);
      if (filters.new_only) params.set('new_only', 'true');

      const response = await api.get(`/export.csv?${params.toString()}`, {
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `matches-${new Date().toISOString().split('T')[0]}.csv`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);

      setState('idle');
      onClose();
    } catch (err: unknown) {
      if (err && typeof err === 'object' && 'response' in err) {
        const axiosError = err as { response?: { status?: number; headers?: Record<string, string>; data?: Blob } };
        const status = axiosError.response?.status;

        if (status === 429) {
          const retryAfter = axiosError.response?.headers?.['retry-after'];
          setState('rate_limited');
          setError({
            code: 'RATE_LIMITED',
            message: 'Export rate limit exceeded. Please wait before trying again.',
            retryAfter: retryAfter ? parseInt(retryAfter, 10) : 60,
          });
          return;
        }

        if (status && status >= 500) {
          setState('error');
          setError({
            code: 'EXPORT_ERROR',
            message: 'Failed to generate export. Please try again later.',
          });
          return;
        }
      }

      setState('error');
      setError({
        code: 'UNKNOWN_ERROR',
        message: 'An unexpected error occurred. Please try again.',
      });
    }
  }, [filters, onClose]);

  const handleClose = useCallback(() => {
    if (state !== 'exporting') {
      setState('idle');
      setError(null);
      onClose();
    }
  }, [state, onClose]);

  const handleRetry = useCallback(() => {
    setState('idle');
    setError(null);
  }, []);

  if (!isOpen) return null;

  const selectedKeywords = keywords.filter((k) => filters.keyword_ids.includes(k.keyword_id));
  const hasFilters = 
    filters.keyword_ids.length > 0 || 
    filters.start_date || 
    filters.end_date || 
    filters.search || 
    filters.new_only;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div 
        className="absolute inset-0 bg-black/50" 
        onClick={handleClose}
        aria-hidden="true"
      />

      <div className="relative bg-white rounded-xl shadow-xl max-w-md w-full mx-4 overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-text-primary">
            Export Matches
          </h2>
        </div>

        <div className="px-6 py-4">
          {state === 'rate_limited' && error && (
            <div className="mb-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
              <div className="flex items-start gap-3">
                <svg className="w-5 h-5 text-yellow-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                <div>
                  <p className="text-sm font-medium text-yellow-800">Rate Limit Exceeded</p>
                  <p className="text-sm text-yellow-700 mt-1">
                    {error.message}
                    {error.retryAfter && (
                      <span> Try again in {error.retryAfter} seconds.</span>
                    )}
                  </p>
                </div>
              </div>
            </div>
          )}

          {state === 'error' && error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <div className="flex items-start gap-3">
                <svg className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
                <div>
                  <p className="text-sm font-medium text-red-800">Export Failed</p>
                  <p className="text-sm text-red-700 mt-1">{error.message}</p>
                </div>
              </div>
            </div>
          )}

          <div className="space-y-3">
            <div>
              <p className="text-sm font-medium text-text-secondary mb-2">
                Filter Summary
              </p>
              {!hasFilters ? (
                <p className="text-sm text-text-muted">
                  Exporting all matches (no filters applied)
                </p>
              ) : (
                <div className="flex flex-wrap gap-2">
                  {selectedKeywords.map((k) => (
                    <FilterChip key={k.keyword_id}>
                      Keyword: {k.value}
                    </FilterChip>
                  ))}
                  {filters.start_date && (
                    <FilterChip>From: {formatDate(filters.start_date)}</FilterChip>
                  )}
                  {filters.end_date && (
                    <FilterChip>To: {formatDate(filters.end_date)}</FilterChip>
                  )}
                  {filters.search && (
                    <FilterChip>Search: "{filters.search}"</FilterChip>
                  )}
                  {filters.new_only && (
                    <FilterChip>New only</FilterChip>
                  )}
                </div>
              )}
            </div>

            <p className="text-xs text-text-muted">
              The export will be downloaded as a CSV file containing all matching certificates.
            </p>
          </div>
        </div>

        <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 flex justify-end gap-3">
          <Button
            variant="secondary"
            onClick={handleClose}
            disabled={state === 'exporting'}
          >
            Cancel
          </Button>
          {(state === 'error' || state === 'rate_limited') ? (
            <Button onClick={handleRetry}>
              Try Again
            </Button>
          ) : (
            <Button
              onClick={handleExport}
              isLoading={state === 'exporting'}
              disabled={state === 'exporting'}
            >
              {state === 'exporting' ? 'Exporting...' : 'Download CSV'}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
