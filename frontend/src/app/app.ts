import { Component, inject, signal } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';

import { AuthService } from './core/auth.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, RouterLink],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {
  private readonly auth = inject(AuthService);

  protected readonly user = this.auth.user;
  protected readonly menuOpen = signal(false);

  constructor() {
    this.auth.loadCurrentUser().subscribe();
  }

  login(): void {
    this.auth.login();
  }

  logout(): void {
    this.auth.logout().subscribe();
  }

  avatarUrl(discordId: string, avatar: string): string {
    return `https://cdn.discordapp.com/avatars/${discordId}/${avatar}.png?size=64`;
  }
}
