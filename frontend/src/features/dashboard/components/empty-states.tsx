import { Button } from '@/components/ui/button';

interface EmptyNoKeywordsProps {
  onNavigateToKeywords: () => void;
}

export function EmptyNoKeywords({ onNavigateToKeywords }: EmptyNoKeywordsProps) {
  return (
    <div className="card text-center py-12">
      <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
      </div>
      <h3 className="text-lg font-medium text-text-primary mb-2">
        No Keywords Configured
      </h3>
      <p className="text-text-muted mb-6 max-w-md mx-auto">
        Add keywords to start monitoring Certificate Transparency logs for potential brand impersonation.
      </p>
      <Button onClick={onNavigateToKeywords}>
        Add Keywords
      </Button>
    </div>
  );
}

export function EmptyNoMatches() {
  return (
    <div className="card text-center py-12">
      <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
        <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 className="text-lg font-medium text-text-primary mb-2">
        No Matches Found
      </h3>
      <p className="text-text-muted max-w-md mx-auto">
        No certificates matching your keywords have been detected yet. 
        The monitor is actively scanning Certificate Transparency logs and will display any matches here.
      </p>
    </div>
  );
}
