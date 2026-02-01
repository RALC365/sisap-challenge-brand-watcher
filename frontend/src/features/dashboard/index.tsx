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
import api from '@/lib/axios';

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
  const [isExporting, setIsExporting] = useState(false);

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

  const handleExport = useCallback(async () => {
    setIsExporting(true);
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
    } catch (error) {
      console.error('Export failed:', error);
    } finally {
      setIsExporting(false);
    }
  }, [filters]);

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
        onExport={handleExport}
        isExporting={isExporting}
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
      </main>
    </div>
  );
}

export default Dashboard;
