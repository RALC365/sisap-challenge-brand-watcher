import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/axios';
import type { Keyword } from '@/lib/schemas';

interface CreateKeywordInput {
  value: string;
}

interface CreateKeywordError {
  code: string;
  message: string;
}

export function useCreateKeyword() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateKeywordInput): Promise<Keyword> => {
      const { data } = await api.post('/keywords', input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['keywords'] });
    },
    onError: (error: unknown) => {
      return error;
    },
  });
}

export function getCreateKeywordError(error: unknown): CreateKeywordError | null {
  if (error && typeof error === 'object' && 'response' in error) {
    const axiosError = error as { response?: { status?: number; data?: { code?: string; message?: string } } };
    if (axiosError.response?.status === 409) {
      return {
        code: axiosError.response.data?.code || 'DUPLICATE_KEYWORD',
        message: axiosError.response.data?.message || 'This keyword already exists',
      };
    }
    if (axiosError.response?.status === 400) {
      return {
        code: axiosError.response.data?.code || 'VALIDATION_ERROR',
        message: axiosError.response.data?.message || 'Invalid keyword',
      };
    }
  }
  return null;
}
