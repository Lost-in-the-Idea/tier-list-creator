import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import {
  CdkDrag,
  CdkDragDrop,
  CdkDropList,
  CdkDropListGroup,
  moveItemInArray,
  transferArrayItem,
} from '@angular/cdk/drag-drop';

import { AuthService } from '../../core/auth.service';
import { TierlistService } from '../../core/tierlist.service';
import { Tierlist, TierlistItem } from '../../core/models';
import { TIERS } from '../../core/tiers';

interface TierBucket {
  label: string;
  color: string;
  items: TierlistItem[];
}

@Component({
  selector: 'app-vote',
  imports: [CdkDropList, CdkDrag, CdkDropListGroup, RouterLink],
  templateUrl: './vote.html',
  styleUrl: './vote.css',
})
export class Vote {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly tierlists = inject(TierlistService);
  private readonly auth = inject(AuthService);

  readonly poolId = 'tier-pool';

  protected readonly user = this.auth.user;
  protected readonly tierlist = signal<Tierlist | null>(null);
  protected readonly loading = signal(true);
  protected readonly error = signal<string | null>(null);
  protected readonly submitting = signal(false);

  protected readonly tiers = signal<TierBucket[]>(
    TIERS.map((t) => ({ label: t.label, color: t.color, items: [] })),
  );
  protected readonly pool = signal<TierlistItem[]>([]);

  protected readonly id = this.route.snapshot.paramMap.get('id') ?? '';

  protected readonly connectedIds = computed(() => [
    this.poolId,
    ...this.tiers().map((t) => t.label),
  ]);
  protected readonly allPlaced = computed(() => this.pool().length === 0);

  constructor() {
    if (!this.id) {
      this.error.set('No tier list specified.');
      this.loading.set(false);
      return;
    }
    this.tierlists.getById(this.id).subscribe({
      next: (tl) => {
        this.loading.set(false);
        const expired = new Date(tl.expires_at).getTime() < Date.now();
        if (tl.has_submitted || expired) {
          this.router.navigate(['/t', this.id, 'results']);
          return;
        }
        this.tierlist.set(tl);
        this.pool.set([...tl.items].sort((a, b) => a.sort_order - b.sort_order));
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(
          err?.status === 404 ? 'Tier list not found.' : 'Failed to load tier list.',
        );
      },
    });
  }

  drop(event: CdkDragDrop<TierlistItem[]>): void {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
    } else {
      transferArrayItem(
        event.previousContainer.data,
        event.container.data,
        event.previousIndex,
        event.currentIndex,
      );
    }
    this.tiers.set([...this.tiers()]);
    this.pool.set([...this.pool()]);
  }

  submit(): void {
    if (!this.allPlaced() || this.submitting()) {
      return;
    }
    if (!this.user()) {
      this.error.set('You need to log in to submit your ranking.');
      return;
    }
    this.error.set(null);
    this.submitting.set(true);

    const rankings = this.tiers().flatMap((tier) =>
      tier.items.map((item) => ({ item_id: item.id, tier: tier.label })),
    );

    this.tierlists.submit(this.id, { rankings }).subscribe({
      next: () => this.router.navigate(['/t', this.id, 'results']),
      error: (err) => {
        this.submitting.set(false);
        if (err?.status === 409) {
          this.router.navigate(['/t', this.id, 'results']);
          return;
        }
        this.error.set(
          err?.status === 401
            ? 'You need to log in to submit your ranking.'
            : err?.error?.error || 'Failed to submit ranking.',
        );
      },
    });
  }
}
