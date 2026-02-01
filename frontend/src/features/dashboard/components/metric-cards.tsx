import { SkeletonCard } from '@/components/feedback/skeleton';
import type { MonitorStatus } from '@/lib/schemas';

interface MetricCardsProps {
  metrics: MonitorStatus['metrics_last_run'] | null | undefined;
  isLoading: boolean;
}

interface MetricCardProps {
  label: string;
  value: number | string;
  variant?: 'default' | 'success' | 'warning' | 'error';
}

const variantStyles = {
  default: 'text-text-primary',
  success: 'text-success',
  warning: 'text-warning',
  error: 'text-error',
};

function MetricCard({ label, value, variant = 'default' }: MetricCardProps) {
  return (
    <div className="card">
      <p className="text-sm text-text-muted mb-1">{label}</p>
      <p className={`text-2xl font-semibold ${variantStyles[variant]}`}>
        {value}
      </p>
    </div>
  );
}

export function MetricCards({ metrics, isLoading }: MetricCardsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <SkeletonCard />
        <SkeletonCard />
        <SkeletonCard />
        <SkeletonCard />
      </div>
    );
  }

  const processedCount = metrics?.processed_count ?? 0;
  const matchCount = metrics?.match_count ?? 0;
  const parseErrorCount = metrics?.parse_error_count ?? 0;
  const durationMs = metrics?.duration_ms ?? 0;

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <MetricCard
        label="Processed"
        value={processedCount.toLocaleString()}
        variant="default"
      />
      <MetricCard
        label="Matches"
        value={matchCount.toLocaleString()}
        variant={matchCount > 0 ? 'success' : 'default'}
      />
      <MetricCard
        label="Parse Errors"
        value={parseErrorCount.toLocaleString()}
        variant={parseErrorCount > 0 ? 'warning' : 'default'}
      />
      <MetricCard
        label="Cycle Duration"
        value={`${durationMs}ms`}
        variant="default"
      />
    </div>
  );
}
