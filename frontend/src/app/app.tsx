import { ToastProvider } from '@/components/feedback/toast';
import { Dashboard } from '@/features/dashboard';
import { KeywordsPage } from '@/features/keywords/pages/keywords-page';

function Router() {
  const path = window.location.pathname;

  if (path === '/keywords') {
    return <KeywordsPage />;
  }

  return <Dashboard />;
}

export function App() {
  return (
    <ToastProvider>
      <Router />
    </ToastProvider>
  );
}
