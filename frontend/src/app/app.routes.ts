import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./pages/home/home').then((m) => m.Home),
  },
  {
    path: 'create',
    loadComponent: () => import('./pages/create/create').then((m) => m.Create),
  },
  {
    path: 't/:id',
    loadComponent: () => import('./pages/vote/vote').then((m) => m.Vote),
  },
  {
    path: 't/:id/results',
    loadComponent: () => import('./pages/results/results').then((m) => m.Results),
  },
  { path: '**', redirectTo: '' },
];
