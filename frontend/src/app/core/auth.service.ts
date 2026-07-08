import { HttpClient } from '@angular/common/http';
import { Injectable, inject, signal } from '@angular/core';
import { catchError, of, tap } from 'rxjs';

import { User } from './models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);

  readonly user = signal<User | null>(null);
  readonly loaded = signal(false);

  /** Fetch the current user from the session cookie (called once on startup). */
  loadCurrentUser() {
    return this.http.get<User>('/api/auth/me').pipe(
      tap((user) => {
        this.user.set(user);
        this.loaded.set(true);
      }),
      catchError(() => {
        this.user.set(null);
        this.loaded.set(true);
        return of(null);
      }),
    );
  }

  /** Full-page redirect into the Discord OAuth flow (must hit the backend directly). */
  login(): void {
    window.location.href = '/api/auth/discord/redirect';
  }

  logout() {
    return this.http.get('/api/auth/logout').pipe(
      tap(() => this.user.set(null)),
      catchError(() => {
        this.user.set(null);
        return of(null);
      }),
    );
  }
}
