import { Table, Pagination } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { SkeletonTable } from '@/components/feedback/skeleton';
import type { Match } from '@/lib/schemas';

function HighlightedText({ text, keyword }: { text: string; keyword: string }) {
  if (!keyword) return <>{text}</>;
  
  const regex = new RegExp(`(${keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
  const parts = text.split(regex);
  
  return (
    <>
      {parts.map((part, i) => 
        regex.test(part) ? (
          <mark key={i} className="bg-yellow-200 text-yellow-900 px-0.5 rounded">{part}</mark>
        ) : (
          <span key={i}>{part}</span>
        )
      )}
    </>
  );
}

function truncateFingerprint(hash: string): string {
  if (hash.length <= 16) return hash;
  return `${hash.slice(0, 8)}...${hash.slice(-8)}`;
}

interface MatchesTableProps {
  matches: Match[];
  total: number;
  page: number;
  limit: number;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  onPageChange: (page: number) => void;
  onSort: (key: string) => void;
  isLoading: boolean;
}

function formatDate(dateString: string | null): string {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function MatchesTable({
  matches,
  total,
  page,
  limit,
  sortBy,
  sortOrder,
  onPageChange,
  onSort,
  isLoading,
}: MatchesTableProps) {
  const totalPages = Math.ceil(total / limit);

  const columns = [
    {
      key: 'keyword_value',
      header: 'Keyword',
      sortable: true,
      render: (match: Match) => (
        <span className="font-medium text-primary">{match.keyword_value}</span>
      ),
    },
    {
      key: 'matched_value',
      header: 'Matched Domain',
      sortable: true,
      className: 'max-w-md',
      render: (match: Match) => (
        <div className="break-words">
          <span className="font-mono text-sm">
            <HighlightedText text={match.matched_value} keyword={match.keyword_value} />
          </span>
          {match.is_new && (
            <Badge variant="success" className="ml-2">New</Badge>
          )}
        </div>
      ),
    },
    {
      key: 'matched_field',
      header: 'Field',
      render: (match: Match) => (
        <Badge variant="info">{match.matched_field.toUpperCase()}</Badge>
      ),
    },
    {
      key: 'certificate_sha256',
      header: 'Fingerprint',
      className: 'max-w-[120px]',
      render: (match: Match) => (
        <span 
          className="font-mono text-xs text-text-muted cursor-help truncate block" 
          title={match.certificate_sha256}
        >
          {truncateFingerprint(match.certificate_sha256)}
        </span>
      ),
    },
    {
      key: 'issuer_org',
      header: 'Issuer',
      className: 'max-w-xs',
      render: (match: Match) => (
        <span className="text-text-muted break-words">
          {match.issuer_org || match.issuer_cn || '-'}
        </span>
      ),
    },
    {
      key: 'first_seen_at',
      header: 'First Seen',
      sortable: true,
      render: (match: Match) => (
        <span className="text-sm text-text-muted whitespace-nowrap">
          {formatDate(match.first_seen_at)}
        </span>
      ),
    },
    {
      key: 'not_after',
      header: 'Expires',
      sortable: true,
      render: (match: Match) => (
        <span className="text-sm text-text-muted whitespace-nowrap">
          {formatDate(match.not_after)}
        </span>
      ),
    },
  ];

  if (isLoading) {
    return (
      <div className="card overflow-hidden">
        <SkeletonTable rows={10} />
      </div>
    );
  }

  return (
    <div className="card overflow-hidden p-0">
      <div className="max-h-[600px] overflow-auto">
        <Table
          columns={columns}
          data={matches}
          sortBy={sortBy}
          sortOrder={sortOrder}
          onSort={onSort}
          keyExtractor={(match) => match.id}
          emptyMessage="No matches found"
        />
      </div>
      {total > 0 && (
        <Pagination
          currentPage={page}
          totalPages={totalPages}
          totalItems={total}
          itemsPerPage={limit}
          onPageChange={onPageChange}
        />
      )}
    </div>
  );
}
