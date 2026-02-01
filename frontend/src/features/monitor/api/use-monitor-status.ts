import { useQuery } from '@tanstack/react-query';
import api from '@/lib/axios';
import { MonitorStatusSchema, type MonitorStatus } from '@/lib/schemas';

export function useMonitorStatus() {
  return useQuery({
    queryKey: ['monitor', 'status'],
    queryFn: async (): Promise<MonitorStatus> => {
      const { data } = await api.get('/monitor/status');
      return MonitorStatusSchema.parse(data);
    },
    refetchInterval: 10000,
  });
}
