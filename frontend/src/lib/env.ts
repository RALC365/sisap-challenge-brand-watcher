export const env = {
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || '/api',
  isDev: import.meta.env.DEV,
  isProd: import.meta.env.PROD,
};
