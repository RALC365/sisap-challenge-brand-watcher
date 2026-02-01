import { StatusBadge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { MonitorStatus } from '@/lib/schemas';

interface AppBarProps {
  status: MonitorStatus | undefined;
  isLoading: boolean;
  onExport: () => void;
}

function formatLastRun(dateString: string | null): string {
  if (!dateString) return 'Never';
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  
  if (diffMins < 1) return 'Just now';
  if (diffMins < 60) return `${diffMins}m ago`;
  
  const diffHours = Math.floor(diffMins / 60);
  if (diffHours < 24) return `${diffHours}h ago`;
  
  return date.toLocaleDateString();
}

export function AppBar({ status, isLoading, onExport }: AppBarProps) {
  return (
    <header className="sticky top-0 z-20 bg-surface-card shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h1 className="text-xl font-semibold text-text-primary">
              Brand Protection Monitor
            </h1>
            {isLoading ? (
              <div className="animate-pulse bg-gray-200 h-6 w-16 rounded-full" />
            ) : status ? (
              <StatusBadge state={status.state} />
            ) : null}
          </div>

          <div className="flex items-center gap-4">
            {status && (
              <span className="text-sm text-text-muted">
                Last run: {formatLastRun(status.last_run_at)}
              </span>
            )}
            <Button
              onClick={onExport}
              variant="primary"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              Export CSV
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
}
