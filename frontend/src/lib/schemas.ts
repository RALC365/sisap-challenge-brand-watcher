import { z } from 'zod';

export const KeywordSchema = z.object({
  keyword_id: z.string().uuid(),
  value: z.string().min(1).max(64),
  normalized_value: z.string(),
  status: z.enum(['active', 'inactive']),
  created_at: z.string(),
});

export const KeywordListResponseSchema = z.object({
  items: z.array(KeywordSchema),
  total: z.number(),
});

export const MatchSchema = z.object({
  id: z.string().uuid(),
  keyword_id: z.string().uuid(),
  keyword_value: z.string(),
  certificate_sha256: z.string(),
  matched_field: z.enum(['cn', 'san', 'both']),
  matched_value: z.string(),
  domain_name: z.string().nullable(),
  issuer_cn: z.string().nullable(),
  issuer_org: z.string().nullable(),
  subject_cn: z.string().nullable(),
  subject_org: z.string().nullable(),
  not_before: z.string().nullable(),
  not_after: z.string().nullable(),
  first_seen_at: z.string(),
  last_seen_at: z.string(),
  is_new: z.boolean(),
  ct_log_index: z.number(),
});

export const MatchListResponseSchema = z.object({
  items: z.array(MatchSchema),
  total: z.number(),
});

export const MonitorStatusSchema = z.object({
  state: z.enum(['idle', 'running', 'error']),
  last_run_at: z.string().nullable(),
  last_success_at: z.string().nullable(),
  last_error_code: z.string().nullable(),
  last_error_message: z.string().nullable(),
  metrics_last_run: z.object({
    processed_count: z.number(),
    match_count: z.number(),
    parse_error_count: z.number(),
    duration_ms: z.number(),
    ct_latency_ms: z.number(),
    db_latency_ms: z.number(),
  }).nullable(),
});

export const CreateKeywordSchema = z.object({
  value: z.string().min(1).max(64),
});

export type Keyword = z.infer<typeof KeywordSchema>;
export type KeywordListResponse = z.infer<typeof KeywordListResponseSchema>;
export type Match = z.infer<typeof MatchSchema>;
export type MatchListResponse = z.infer<typeof MatchListResponseSchema>;
export type MonitorStatus = z.infer<typeof MonitorStatusSchema>;
export type CreateKeyword = z.infer<typeof CreateKeywordSchema>;
