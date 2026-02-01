import { useQuery } from '@tanstack/react-query';
import api from '@/lib/axios';
import { KeywordListResponseSchema, type KeywordListResponse } from '@/lib/schemas';

export function useKeywords() {
  return useQuery({
    queryKey: ['keywords'],
    queryFn: async (): Promise<KeywordListResponse> => {
      const { data } = await api.get('/keywords');
      return KeywordListResponseSchema.parse(data);
    },
  });
}
