import { useState, useCallback, type FormEvent } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useCreateKeyword, getCreateKeywordError } from '../api/use-create-keyword';

interface KeywordFormProps {
  onSuccess?: () => void;
}

export function KeywordForm({ onSuccess }: KeywordFormProps) {
  const [value, setValue] = useState('');
  const [error, setError] = useState<string | null>(null);
  const createKeyword = useCreateKeyword();

  const validate = useCallback((input: string): string | null => {
    const trimmed = input.trim();
    if (!trimmed) {
      return 'Keyword is required';
    }
    if (trimmed.length < 2) {
      return 'Keyword must be at least 2 characters';
    }
    if (trimmed.length > 64) {
      return 'Keyword must be at most 64 characters';
    }
    return null;
  }, []);

  const handleSubmit = useCallback(async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    const validationError = validate(value);
    if (validationError) {
      setError(validationError);
      return;
    }

    try {
      await createKeyword.mutateAsync({ value: value.trim() });
      setValue('');
      onSuccess?.();
    } catch (err) {
      const apiError = getCreateKeywordError(err);
      if (apiError?.code === 'DUPLICATE_KEYWORD') {
        setError('This keyword already exists');
      } else if (apiError) {
        setError(apiError.message);
      } else {
        setError('Failed to create keyword');
      }
    }
  }, [value, validate, createKeyword, onSuccess]);

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    if (error) {
      setError(null);
    }
  }, [error]);

  return (
    <form onSubmit={handleSubmit} className="flex gap-3 items-start">
      <div className="flex-1">
        <Input
          value={value}
          onChange={handleChange}
          placeholder="Enter keyword to monitor..."
          error={error || undefined}
          disabled={createKeyword.isPending}
        />
      </div>
      <Button
        type="submit"
        isLoading={createKeyword.isPending}
        disabled={!value.trim()}
      >
        Add Keyword
      </Button>
    </form>
  );
}
