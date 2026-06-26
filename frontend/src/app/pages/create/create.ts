import { Component, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from '../../core/auth.service';
import { TierlistService } from '../../core/tierlist.service';

interface ItemDraft {
  name: string;
  image_url: string;
}

@Component({
  selector: 'app-create',
  imports: [],
  templateUrl: './create.html',
})
export class Create {
  private readonly tierlists = inject(TierlistService);
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly user = this.auth.user;

  protected readonly title = signal('');
  protected readonly description = signal('');
  protected readonly expiry = signal(this.defaultExpiry());
  protected readonly items = signal<ItemDraft[]>([
    { name: '', image_url: '' },
    { name: '', image_url: '' },
  ]);

  protected readonly submitting = signal(false);
  protected readonly error = signal<string | null>(null);

  protected readonly filledItems = computed(
    () => this.items().filter((i) => i.name.trim()).length,
  );
  protected readonly canSubmit = computed(
    () =>
      !!this.title().trim() &&
      !!this.expiry() &&
      this.filledItems() >= 2 &&
      this.filledItems() <= 15,
  );

  addItem(): void {
    if (this.items().length < 15) {
      this.items.update((items) => [...items, { name: '', image_url: '' }]);
    }
  }

  removeItem(index: number): void {
    this.items.update((items) => items.filter((_, i) => i !== index));
  }

  updateItem(index: number, field: keyof ItemDraft, value: string): void {
    this.items.update((items) =>
      items.map((item, i) => (i === index ? { ...item, [field]: value } : item)),
    );
  }

  submit(): void {
    if (!this.canSubmit() || this.submitting()) {
      return;
    }
    this.error.set(null);
    this.submitting.set(true);

    const tierlist_items = this.items()
      .map((item, index) => ({
        name: item.name.trim(),
        image_url: item.image_url.trim(),
        sort_order: index,
      }))
      .filter((item) => item.name);

    this.tierlists
      .create({
        title: this.title().trim(),
        description: this.description().trim(),
        expiry_time: new Date(this.expiry()).toISOString(),
        tierlist_items,
      })
      .subscribe({
        next: (res) => this.router.navigate(['/t', res.id]),
        error: (err) => {
          this.submitting.set(false);
          this.error.set(
            err?.status === 401
              ? 'You need to log in to create a tier list.'
              : err?.error?.error || 'Something went wrong creating the tier list.',
          );
        },
      });
  }

  private defaultExpiry(): string {
    const d = new Date();
    d.setDate(d.getDate() + 7);
    d.setMinutes(d.getMinutes() - d.getTimezoneOffset());
    return d.toISOString().slice(0, 16);
  }
}
