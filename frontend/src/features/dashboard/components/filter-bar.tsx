import { useState, useEffect, useCallback } from 'react';
import { Input } from '@/components/ui/input';
import { Toggle } from '@/components/ui/toggle';
import { DateRange } from '@/components/ui/date-range';
import type { Keyword } from '@/lib/schemas';
import { useDebounce } from '../hooks/use-debounce';

export interface FilterState {
  keyword_ids: string[];
  start_date: string;
  end_date: string;
  search: string;
  new_only: boolean;
}

interface FilterBarProps {
  keywords: Keyword[];
  filters: FilterState;
  onFiltersChange: (filters: FilterState) => void;
  isLoading?: boolean;
}

export function FilterBar({ keywords, filters, onFiltersChange, isLoading }: FilterBarProps) {
  const [searchInput, setSearchInput] = useState(filters.search);
  const debouncedSearch = useDebounce(searchInput, 300);

  useEffect(() => {
    if (debouncedSearch !== filters.search) {
      onFiltersChange({ ...filters, search: debouncedSearch });
    }
  }, [debouncedSearch, filters, onFiltersChange]);

  const handleKeywordToggle = useCallback((keywordId: string) => {
    const newIds = filters.keyword_ids.includes(keywordId)
      ? filters.keyword_ids.filter((id) => id !== keywordId)
      : [...filters.keyword_ids, keywordId];
    onFiltersChange({ ...filters, keyword_ids: newIds });
  }, [filters, onFiltersChange]);

  const handleDateChange = useCallback((start: string, end: string) => {
    onFiltersChange({ ...filters, start_date: start, end_date: end });
  }, [filters, onFiltersChange]);

  const handleNewOnlyChange = useCallback((checked: boolean) => {
    onFiltersChange({ ...filters, new_only: checked });
  }, [filters, onFiltersChange]);

  const clearFilters = useCallback(() => {
    setSearchInput('');
    onFiltersChange({
      keyword_ids: [],
      start_date: '',
      end_date: '',
      search: '',
      new_only: false,
    });
  }, [onFiltersChange]);

  const hasActiveFilters = filters.keyword_ids.length > 0 ||
    filters.start_date ||
    filters.end_date ||
    filters.search ||
    filters.new_only;

  return (
    <div className="card mb-6">
      <div className="flex flex-col xl:flex-row gap-4">
        <div className="flex-1 min-w-0">
          <label className="label">Keywords</label>
          <div className="flex flex-wrap gap-2 min-h-[38px] p-2 border border-gray-200 rounded-md bg-white">
            {keywords.length === 0 ? (
              <span className="text-sm text-text-muted">No keywords configured</span>
            ) : (
              keywords.map((keyword) => (
                <button
                  key={keyword.keyword_id}
                  onClick={() => handleKeywordToggle(keyword.keyword_id)}
                  disabled={isLoading}
                  className={`px-2 py-1 text-xs rounded-full transition-colors ${
                    filters.keyword_ids.includes(keyword.keyword_id)
                      ? 'bg-primary text-white'
                      : 'bg-gray-100 text-text-primary hover:bg-gray-200'
                  }`}
                >
                  {keyword.value}
                </button>
              ))
            )}
          </div>
        </div>

        <div className="shrink-0">
          <DateRange
            label="Date Range"
            startDate={filters.start_date}
            endDate={filters.end_date}
            onStartDateChange={(date) => handleDateChange(date, filters.end_date)}
            onEndDateChange={(date) => handleDateChange(filters.start_date, date)}
          />
        </div>

        <div className="w-full xl:w-40 shrink-0">
          <Input
            label="Search"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Search..."
          />
        </div>

        <div className="flex items-end gap-3 shrink-0">
          <div className="pb-2">
            <Toggle
              label="New only"
              checked={filters.new_only}
              onChange={handleNewOnlyChange}
            />
          </div>

          {hasActiveFilters && (
            <button
              onClick={clearFilters}
              className="text-sm text-primary hover:text-blue-700 pb-2 whitespace-nowrap"
            >
              Clear filters
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
