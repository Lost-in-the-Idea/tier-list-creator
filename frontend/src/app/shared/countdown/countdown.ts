import { Component, DestroyRef, computed, effect, inject, input, output, signal } from '@angular/core';

import { formatCountdown } from './countdown-format';

@Component({
  selector: 'app-countdown',
  imports: [],
  templateUrl: './countdown.html',
})
export class Countdown {
  readonly expiresAt = input.required<string | Date>();
  readonly expired = output<void>();

  private readonly destroyRef = inject(DestroyRef);
  private readonly now = signal(Date.now());
  private hasEmittedExpired = false;

  private readonly expiresAtMs = computed(() => {
    const v = this.expiresAt();
    return (v instanceof Date ? v : new Date(v)).getTime();
  });

  protected readonly remainingMs = computed(() => Math.max(0, this.expiresAtMs() - this.now()));
  protected readonly isExpired = computed(() => this.remainingMs() <= 0);
  protected readonly display = computed(() => formatCountdown(this.remainingMs()));
  protected readonly isUrgent = computed(
    () => this.remainingMs() > 0 && this.remainingMs() < 5 * 60 * 1000,
  );

  constructor() {
    const intervalId = setInterval(() => this.now.set(Date.now()), 1000);
    this.destroyRef.onDestroy(() => clearInterval(intervalId));

    effect(() => {
      if (this.isExpired() && !this.hasEmittedExpired) {
        this.hasEmittedExpired = true;
        this.expired.emit();
      }
    });
  }
}
