import { DecimalPipe } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import { TierlistService } from '../../core/tierlist.service';
import { TierResult, TierlistResultResponse } from '../../core/models';
import { TIERS, TierDef } from '../../core/tiers';

interface ResultBucket extends TierDef {
  items: TierResult[];
}

type ResultsView = 'board' | 'details';

@Component({
  selector: 'app-results',
  imports: [RouterLink, DecimalPipe],
  templateUrl: './results.html',
  styleUrl: './results.css',
})
export class Results {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly tierlists = inject(TierlistService);

  protected readonly tierDefs: ReadonlyArray<TierDef> = TIERS;
  protected readonly data = signal<TierlistResultResponse | null>(null);
  protected readonly loading = signal(true);
  protected readonly error = signal<string | null>(null);
  protected readonly id = this.route.snapshot.paramMap.get('id') ?? '';
  protected readonly view = signal<ResultsView>(
    this.route.snapshot.queryParamMap.get('view') === 'details' ? 'details' : 'board',
  );

  setView(view: ResultsView): void {
    this.view.set(view);
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { view },
      queryParamsHandling: 'merge',
      replaceUrl: true,
    });
  }

  /** Aggregated results laid out as a read-only tier board (items in their top tier). */
  protected readonly board = computed(() => {
    const buckets: ResultBucket[] = this.tierDefs.map((t) => ({ ...t, items: [] }));
    const unranked: TierResult[] = [];
    const data = this.data();
    if (data) {
      for (const result of data.results) {
        const bucket = buckets.find((b) => b.label === result.top_tier);
        if (bucket) {
          bucket.items.push(result);
        } else {
          unranked.push(result);
        }
      }
    }
    return { buckets, unranked };
  });

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
