import { useQuery, keepPreviousData } from '@tanstack/react-query';
import api from '@/lib/axios';
import { MatchListResponseSchema, type MatchListResponse } from '@/lib/schemas';
import { buildQueryString } from '@/lib/url';

export interface MatchFilters {
  keyword_ids?: string[];
  start_date?: string;
  end_date?: string;
  search?: string;
  new_only?: boolean;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export function useMatches(filters: MatchFilters = {}) {
  const { keyword_ids, ...rest } = filters;
  
  return useQuery({
    queryKey: ['matches', filters],
    queryFn: async (): Promise<MatchListResponse> => {
      const params: Record<string, string | number | boolean | undefined> = {
        ...rest,
        keyword_ids: keyword_ids?.join(','),
      };
      const queryString = buildQueryString(params);
      const { data } = await api.get(`/matches${queryString}`);
      return MatchListResponseSchema.parse(data);
    },
    placeholderData: keepPreviousData,
  });
}
