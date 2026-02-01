import { useState, useCallback } from 'react';
import { useKeywords } from '../api/use-keywords';
import { KeywordForm } from '../components/keyword-form';
import { KeywordTable } from '../components/keyword-table';
import { ErrorBanner } from '@/components/feedback/error-banner';
import { Button } from '@/components/ui/button';

export function KeywordsPage() {
  const [bannerError, setBannerError] = useState<string | null>(null);
  const { data, isLoading, error, refetch } = useKeywords();

  const handleFormSuccess = useCallback(() => {
    refetch();
  }, [refetch]);

  const handleServerError = useCallback((message: string) => {
    setBannerError(message);
  }, []);

  const handleDismissError = useCallback(() => {
    setBannerError(null);
  }, []);

  const handleNavigateBack = useCallback(() => {
    window.location.href = '/';
  }, []);

  return (
    <div className="min-h-screen bg-surface-page">
      <header className="sticky top-0 z-20 bg-surface-card shadow-sm border-b border-gray-200">
        <div className="max-w-3xl mx-auto px-4 py-4">
          <div className="flex items-center gap-4">
            <button
              onClick={handleNavigateBack}
              className="text-text-muted hover:text-text-primary"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
            </button>
            <h1 className="text-xl font-semibold text-text-primary">
              Keyword Management
            </h1>
          </div>
        </div>
      </header>

      <main className="max-w-3xl mx-auto px-4 py-6">
        {bannerError && (
          <ErrorBanner
            errorMessage={bannerError}
            onDismiss={handleDismissError}
          />
        )}

        {error && (
          <div className="card mb-6">
            <div className="text-center py-8">
              <p className="text-error mb-4">Failed to load keywords</p>
              <Button onClick={() => refetch()}>Retry</Button>
            </div>
          </div>
        )}

        <div className="card mb-6">
          <h2 className="text-lg font-medium text-text-primary mb-4">
            Add New Keyword
          </h2>
          <KeywordForm onSuccess={handleFormSuccess} />
          <p className="text-xs text-text-muted mt-2">
            Keywords are matched case-insensitively against certificate Common Names (CN) and Subject Alternative Names (SAN).
          </p>
        </div>

        <div className="card">
          <h2 className="text-lg font-medium text-text-primary mb-4">
            Active Keywords ({data?.items.length ?? 0})
          </h2>
          <KeywordTable
            keywords={data?.items ?? []}
            isLoading={isLoading}
            onError={handleServerError}
          />
        </div>
      </main>
    </div>
  );
}

export default KeywordsPage;
