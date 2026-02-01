import { type ReactNode } from 'react';

type BadgeVariant = 'default' | 'success' | 'warning' | 'error' | 'info';

interface BadgeProps {
  variant?: BadgeVariant;
  children: ReactNode;
  className?: string;
}

const variantStyles: Record<BadgeVariant, string> = {
  default: 'bg-gray-100 text-gray-800',
  success: 'bg-green-100 text-green-800',
  warning: 'bg-yellow-100 text-yellow-800',
  error: 'bg-red-100 text-red-800',
  info: 'bg-blue-100 text-blue-800',
};

export function Badge({ variant = 'default', children, className = '' }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${variantStyles[variant]} ${className}`}
    >
      {children}
    </span>
  );
}

interface StatusBadgeProps {
  state: 'idle' | 'running' | 'error';
}

export function StatusBadge({ state }: StatusBadgeProps) {
  const config = {
    idle: { variant: 'success' as const, label: 'Idle', dot: 'bg-green-500' },
    running: { variant: 'info' as const, label: 'Running', dot: 'bg-blue-500 animate-pulse' },
    error: { variant: 'error' as const, label: 'Error', dot: 'bg-red-500' },
  };

  const { variant, label, dot } = config[state];

  return (
    <Badge variant={variant} className="gap-1.5">
      <span className={`w-2 h-2 rounded-full ${dot}`} />
      {label}
    </Badge>
  );
}
