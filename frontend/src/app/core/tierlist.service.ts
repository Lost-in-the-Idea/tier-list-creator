import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';

import {
  CreateTierlistRequest,
  CreateTierlistResponse,
  SubmitRankingRequest,
  Tierlist,
  TierlistResultResponse,
} from './models';

@Injectable({ providedIn: 'root' })
export class TierlistService {
  private readonly http = inject(HttpClient);
  private readonly base = '/api/tierlists';

  create(body: CreateTierlistRequest) {
    return this.http.post<CreateTierlistResponse>(`${this.base}/`, body);
  }

  getById(id: string) {
    return this.http.get<Tierlist>(`${this.base}/${id}`);
  }

  getResults(id: string) {
    return this.http.get<TierlistResultResponse>(`${this.base}/${id}/results`);
  }

  submit(id: string, body: SubmitRankingRequest) {
    return this.http.post<{ message: string }>(`${this.base}/${id}/submit`, body);
  }
}
