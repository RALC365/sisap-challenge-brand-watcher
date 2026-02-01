import { useState, useEffect } from 'react';
import { StatusBadge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { MonitorStatus } from '@/lib/schemas';

interface AppBarProps {
  status: MonitorStatus | undefined;
  isLoading: boolean;
  onExport: () => void;
}

function useCountdown(lastRunAt: string | null, pollIntervalSeconds: number) {
  const [secondsRemaining, setSecondsRemaining] = useState<number | null>(null);

  useEffect(() => {
    if (!lastRunAt || pollIntervalSeconds <= 0) {
      setSecondsRemaining(null);
      return;
    }

    const calculateRemaining = () => {
      const lastRun = new Date(lastRunAt).getTime();
      const nextRun = lastRun + (pollIntervalSeconds * 1000);
      const now = Date.now();
      const remaining = Math.max(0, Math.ceil((nextRun - now) / 1000));
      return remaining;
    };

    setSecondsRemaining(calculateRemaining());

    const interval = setInterval(() => {
      const remaining = calculateRemaining();
      setSecondsRemaining(remaining);
    }, 1000);

    return () => clearInterval(interval);
  }, [lastRunAt, pollIntervalSeconds]);

  return secondsRemaining;
}

export function AppBar({ status, isLoading, onExport }: AppBarProps) {
  const secondsRemaining = useCountdown(
    status?.last_run_at ?? null, 
    status?.poll_interval_seconds ?? 60
  );

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
              <div className="flex items-center gap-3">
                {secondsRemaining !== null && secondsRemaining > 0 && (
                  <span className="text-sm font-mono bg-gray-100 px-2 py-1 rounded text-text-secondary">
                    Next Run: {secondsRemaining}s
                  </span>
                )}
                {secondsRemaining === 0 && status.state !== 'running' && (
                  <span className="text-sm font-mono bg-blue-100 px-2 py-1 rounded text-blue-700 animate-pulse">
                    Polling...
                  </span>
                )}
              </div>
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
