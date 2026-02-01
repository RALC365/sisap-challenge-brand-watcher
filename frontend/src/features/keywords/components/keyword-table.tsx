import { useState, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { useToast } from '@/components/feedback/toast';
import { useDeleteKeyword, getDeleteKeywordError } from '../api/use-delete-keyword';
import type { Keyword } from '@/lib/schemas';

interface KeywordTableProps {
  keywords: Keyword[];
  isLoading: boolean;
  onError: (message: string) => void;
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
}

export function KeywordTable({ keywords, isLoading, onError }: KeywordTableProps) {
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const deleteKeyword = useDeleteKeyword();
  const { addToast } = useToast();

  const handleDelete = useCallback(async (keyword: Keyword) => {
    if (deletingId) return;
    
    setDeletingId(keyword.keyword_id);
    try {
      await deleteKeyword.mutateAsync(keyword.keyword_id);
      addToast('success', `Keyword "${keyword.value}" deleted`);
    } catch (err) {
      const apiError = getDeleteKeywordError(err);
      if (apiError?.status === 404) {
        addToast('error', 'Keyword not found');
      } else if (apiError?.status === 500 || (apiError && apiError.status >= 500)) {
        onError(apiError.message || 'Server error occurred');
      } else {
        addToast('error', apiError?.message || 'Failed to delete keyword');
      }
    } finally {
      setDeletingId(null);
    }
  }, [deletingId, deleteKeyword, addToast, onError]);

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((i) => (
          <div key={i} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg animate-pulse">
            <div className="h-4 w-32 bg-gray-200 rounded" />
            <div className="h-4 w-24 bg-gray-200 rounded" />
          </div>
        ))}
      </div>
    );
  }

  if (keywords.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
          <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-text-primary mb-2">No Keywords Yet</h3>
        <p className="text-text-muted">Add your first keyword to start monitoring Certificate Transparency logs.</p>
      </div>
    );
  }

  return (
    <div className="divide-y divide-gray-200">
      {keywords.map((keyword) => (
        <div
          key={keyword.keyword_id}
          className="flex items-center justify-between py-4 first:pt-0 last:pb-0"
        >
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-text-primary truncate">
              {keyword.value}
            </p>
            <p className="text-xs text-text-muted">
              Added {formatDate(keyword.created_at)}
            </p>
          </div>
          <div className="flex items-center gap-2 ml-4">
            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
              keyword.status === 'active' 
                ? 'bg-green-100 text-green-800' 
                : 'bg-gray-100 text-gray-800'
            }`}>
              {keyword.status}
            </span>
            <Button
              variant="danger"
              size="sm"
              onClick={() => handleDelete(keyword)}
              isLoading={deletingId === keyword.keyword_id}
              disabled={deletingId !== null}
            >
              Delete
            </Button>
          </div>
        </div>
      ))}
    </div>
  );
}
