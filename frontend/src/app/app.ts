import { Component, inject, signal } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';

import { AuthService } from './core/auth.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, RouterLink],
  templateUrl: './app.html',
})
export class App {
  private readonly auth = inject(AuthService);

  protected readonly user = this.auth.user;
  protected readonly menuOpen = signal(false);
  protected readonly theme = signal<'light' | 'dark'>(this.initialTheme());

  constructor() {
    this.auth.loadCurrentUser().subscribe();
  }

  login(): void {
    this.auth.login();
  }

  logout(): void {
    this.auth.logout().subscribe();
  }

  toggleTheme(): void {
    const next = this.theme() === 'dark' ? 'light' : 'dark';
    this.theme.set(next);
    document.documentElement.setAttribute('data-theme', next);
    try {
      localStorage.setItem('theme', next);
    } catch {
      // ignore storage failures (e.g. private mode)
    }
  }

  private initialTheme(): 'light' | 'dark' {
    try {
      const stored = localStorage.getItem('theme');
      if (stored === 'light' || stored === 'dark') {
        return stored;
      }
    } catch {
      // ignore
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  avatarUrl(discordId: string, avatar: string): string {
    return `https://cdn.discordapp.com/avatars/${discordId}/${avatar}.png?size=64`;
  }
}
