export function App() {
  return (
    <div className="min-h-screen bg-surface-page">
      <header className="bg-surface-card shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-h1">Brand Protection Monitor</h1>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6">
        <div className="card">
          <h2 className="text-h2 mb-4">Welcome</h2>
          <p className="text-text-muted mb-4">
            Monitor Certificate Transparency logs for potential brand impersonation and suspicious domain registrations.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            <div className="card">
              <h3 className="text-h3 text-primary">Keywords</h3>
              <p className="text-body-sm text-text-muted mt-2">
                Configure brand keywords to monitor
              </p>
            </div>

            <div className="card">
              <h3 className="text-h3 text-success">Matches</h3>
              <p className="text-body-sm text-text-muted mt-2">
                View detected certificate matches
              </p>
            </div>

            <div className="card">
              <h3 className="text-h3 text-warning">Export</h3>
              <p className="text-body-sm text-text-muted mt-2">
                Download match data as CSV
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
