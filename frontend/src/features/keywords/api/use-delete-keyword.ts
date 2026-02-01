import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/axios';

interface DeleteKeywordError {
  status: number;
  code: string;
  message: string;
}

export function useDeleteKeyword() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (keywordId: string): Promise<void> => {
      await api.delete(`/keywords/${keywordId}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['keywords'] });
      queryClient.invalidateQueries({ queryKey: ['matches'] });
    },
  });
}

export function getDeleteKeywordError(error: unknown): DeleteKeywordError | null {
  if (error && typeof error === 'object' && 'response' in error) {
    const axiosError = error as { response?: { status?: number; data?: { code?: string; message?: string } } };
    const status = axiosError.response?.status || 500;
    
    if (status === 404) {
      return {
        status: 404,
        code: 'NOT_FOUND',
        message: 'Keyword not found',
      };
    }
    
    return {
      status,
      code: axiosError.response?.data?.code || 'SERVER_ERROR',
      message: axiosError.response?.data?.message || 'Failed to delete keyword',
    };
  }
  return null;
}
