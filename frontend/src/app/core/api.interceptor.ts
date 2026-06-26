import { HttpInterceptorFn } from '@angular/common/http';

/** Ensure every API request sends the session cookie. */
export const credentialsInterceptor: HttpInterceptorFn = (req, next) =>
  next(req.clone({ withCredentials: true }));
