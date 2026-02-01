import { useState, useCallback } from 'react';
import { useMonitorStatus } from '@/features/monitor/api/use-monitor-status';
import { useKeywords } from '@/features/keywords/api/use-keywords';
import { useMatches, type MatchFilters } from '@/features/matches/api/use-matches';
import { ErrorBanner } from '@/components/feedback/error-banner';
import { AppBar } from './components/app-bar';
import { MetricCards } from './components/metric-cards';
import { FilterBar, type FilterState } from './components/filter-bar';
import { MatchesTable } from './components/matches-table';
import { EmptyNoKeywords, EmptyNoMatches } from './components/empty-states';
import { ExportModal } from '@/features/export/components/export-modal';

const DEFAULT_FILTERS: FilterState = {
  keyword_ids: [],
  start_date: '',
  end_date: '',
  search: '',
  new_only: false,
};

const ITEMS_PER_PAGE = 20;

export function Dashboard() {
  const [filters, setFilters] = useState<FilterState>(DEFAULT_FILTERS);
  const [page, setPage] = useState(1);
  const [sortBy, setSortBy] = useState('first_seen_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [isExportModalOpen, setIsExportModalOpen] = useState(false);

  const { data: status, isLoading: statusLoading } = useMonitorStatus();
  const { data: keywordsData, isLoading: keywordsLoading } = useKeywords();

  const matchFilters: MatchFilters = {
    keyword_ids: filters.keyword_ids.length > 0 ? filters.keyword_ids : undefined,
    start_date: filters.start_date || undefined,
    end_date: filters.end_date || undefined,
    search: filters.search || undefined,
    new_only: filters.new_only || undefined,
    page,
    limit: ITEMS_PER_PAGE,
    sort_by: sortBy,
    sort_order: sortOrder,
  };

  const { data: matchesData, isLoading: matchesLoading } = useMatches(matchFilters);

  const handleFiltersChange = useCallback((newFilters: FilterState) => {
    setFilters(newFilters);
    setPage(1);
  }, []);

  const handleSort = useCallback((key: string) => {
    if (key === sortBy) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(key);
      setSortOrder('desc');
    }
    setPage(1);
  }, [sortBy]);

  const handleOpenExportModal = useCallback(() => {
    setIsExportModalOpen(true);
  }, []);

  const handleCloseExportModal = useCallback(() => {
    setIsExportModalOpen(false);
  }, []);

  const handleNavigateToKeywords = useCallback(() => {
    window.location.pathname = '/keywords';
  }, []);

  const keywords = keywordsData?.items ?? [];
  const matches = matchesData?.items ?? [];
  const totalMatches = matchesData?.total ?? 0;
  const hasKeywords = keywords.length > 0;
  const hasMatches = totalMatches > 0;
  const isInitialLoading = statusLoading || keywordsLoading;

  return (
    <div className="min-h-screen bg-surface-page">
      <AppBar
        status={status}
        isLoading={statusLoading}
        onExport={handleOpenExportModal}
      />

      <ExportModal
        isOpen={isExportModalOpen}
        onClose={handleCloseExportModal}
        filters={filters}
        keywords={keywords}
      />

      <main className="max-w-7xl mx-auto px-4 py-6">
        {status?.state === 'error' && (
          <ErrorBanner
            errorCode={status.last_error_code}
            errorMessage={status.last_error_message}
          />
        )}

        <MetricCards
          metrics={status?.metrics_last_run}
          isLoading={statusLoading}
        />

        {!isInitialLoading && !hasKeywords ? (
          <EmptyNoKeywords onNavigateToKeywords={handleNavigateToKeywords} />
        ) : (
          <>
            <FilterBar
              keywords={keywords}
              filters={filters}
              onFiltersChange={handleFiltersChange}
              isLoading={keywordsLoading}
            />

            {!matchesLoading && !hasMatches ? (
              <EmptyNoMatches />
            ) : (
              <MatchesTable
                matches={matches}
                total={totalMatches}
                page={page}
                limit={ITEMS_PER_PAGE}
                sortBy={sortBy}
                sortOrder={sortOrder}
                onPageChange={setPage}
                onSort={handleSort}
                isLoading={matchesLoading}
              />
            )}
          </>
        )}

        <div className="mt-8 p-6 bg-white rounded-xl border border-gray-200">
          <h2 className="text-lg font-semibold text-text-primary mb-4">How It Works</h2>
          <div className="grid md:grid-cols-2 gap-6 text-sm text-text-secondary">
            <div>
              <h3 className="font-medium text-text-primary mb-2">1. Certificate Monitoring</h3>
              <p>This system connects to a public Certificate Transparency (CT) log and polls for new SSL/TLS certificates every 60 seconds. It fetches the latest batch of certificates from the log.</p>
            </div>
            <div>
              <h3 className="font-medium text-text-primary mb-2">2. Keyword Matching</h3>
              <p>Each certificate's domain names (Common Name and Subject Alternative Names) are checked against your configured keywords. Matches are highlighted and stored for review.</p>
            </div>
            <div>
              <h3 className="font-medium text-text-primary mb-2">3. Brand Protection</h3>
              <p>When a certificate contains one of your keywords (e.g., your brand name), it may indicate phishing attempts or domain abuse. Review matches to identify potential threats.</p>
            </div>
            <div>
              <h3 className="font-medium text-text-primary mb-2">4. Export & Analysis</h3>
              <p>Use the filters to narrow down results by keyword, date range, or search terms. Export matching certificates to CSV for further analysis or reporting.</p>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-100 text-xs text-text-muted">
            <strong>Tip:</strong> Click on a keyword chip to filter matches. The "Next" timer shows when the next polling cycle will run. New matches are marked with a green "New" badge.
          </div>
        </div>
      </main>
    </div>
  );
}

export default Dashboard;
