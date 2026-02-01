import { create } from 'zustand';

interface MonitorState {
  isPolling: boolean;
  lastError: string | null;
  setPolling: (polling: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useMonitorStore = create<MonitorState>((set) => ({
  isPolling: false,
  lastError: null,
  setPolling: (polling) => set({ isPolling: polling }),
  setError: (error) => set({ lastError: error }),
  clearError: () => set({ lastError: null }),
}));

interface UIState {
  sidebarOpen: boolean;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: false,
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
}));
