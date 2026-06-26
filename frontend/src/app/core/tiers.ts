/** The tier set the backend accepts (validation: oneof=S A B C D F). */
export interface TierDef {
  label: string;
  color: string;
}

export const TIERS: ReadonlyArray<TierDef> = [
  { label: 'S', color: '#ff7f7f' },
  { label: 'A', color: '#ffbf7f' },
  { label: 'B', color: '#ffdf7f' },
  { label: 'C', color: '#bfff7f' },
  { label: 'D', color: '#7fdfff' },
  { label: 'F', color: '#bf9fff' },
];
