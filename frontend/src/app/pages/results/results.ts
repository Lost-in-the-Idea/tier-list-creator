import { DecimalPipe } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';

import { TierlistService } from '../../core/tierlist.service';
import { TierlistResultResponse } from '../../core/models';
import { TIERS, TierDef } from '../../core/tiers';

@Component({
  selector: 'app-results',
  imports: [RouterLink, DecimalPipe],
  templateUrl: './results.html',
})
export class Results {
  private readonly route = inject(ActivatedRoute);
  private readonly tierlists = inject(TierlistService);

  protected readonly tierDefs: ReadonlyArray<TierDef> = TIERS;
  protected readonly data = signal<TierlistResultResponse | null>(null);
  protected readonly loading = signal(true);
  protected readonly error = signal<string | null>(null);
  protected readonly id = this.route.snapshot.paramMap.get('id') ?? '';

  constructor() {
    this.tierlists.getResults(this.id).subscribe({
      next: (res) => {
        this.data.set(res);
        this.loading.set(false);
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(
          err?.status === 404 ? 'Tier list not found.' : 'Failed to load results.',
        );
      },
    });
  }

  tierColor(label: string): string {
    return this.tierDefs.find((t) => t.label === label)?.color ?? '#888';
  }
}
