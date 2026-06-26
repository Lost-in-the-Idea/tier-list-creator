// Dev-server proxy so the browser talks to the API same-origin (/api -> backend).
// This keeps the session cookie working without CORS. In Docker the frontend
// container reaches the backend via the service name; locally it's localhost.
const target = process.env.API_PROXY_TARGET || 'http://localhost:8080';

module.exports = {
  '/api': {
    target,
    secure: false,
    changeOrigin: false,
    // Keep cookies scoped to localhost so the session cookie set during the
    // Discord OAuth flow (on :8080) is shared with the dev server (:4200).
    cookieDomainRewrite: 'localhost',
  },
};
