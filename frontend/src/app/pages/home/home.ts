import { Component, inject, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';

import { AuthService } from '../../core/auth.service';

@Component({
  selector: 'app-home',
  imports: [RouterLink],
  templateUrl: './home.html',
})
export class Home {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly user = this.auth.user;
  protected readonly openId = signal('');

  login(): void {
    this.auth.login();
  }

  open(): void {
    const id = this.openId().trim();
    if (id) {
      this.router.navigate(['/t', id]);
    }
  }

  /** Accept a pasted full URL or a raw id. */
  onOpenInput(value: string): void {
    const match = value.match(/t\/([^/?#]+)/);
    this.openId.set(match ? match[1] : value);
  }
}
