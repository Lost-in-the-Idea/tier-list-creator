export interface User {
  id: string;
  discord_id: string;
  username: string;
  avatar: string;
}

export interface TierlistItem {
  id: string;
  name: string;
  image_url: string;
  sort_order: number;
}

export interface Tierlist {
  id: string;
  share_code: string;
  title: string;
  description: string;
  expires_at: string;
  creator: User;
  items: TierlistItem[];
  has_submitted: boolean;
}

export interface CreateItemRequest {
  name: string;
  image_url?: string;
  sort_order: number;
}

export interface CreateTierlistRequest {
  title: string;
  description?: string;
  expiry_time: string;
  tierlist_items: CreateItemRequest[];
}

export interface CreateTierlistResponse {
  id: string;
  share_code: string;
  expires_at: string;
}

export interface RankingRequest {
  item_id: string;
  tier: string;
}

export interface SubmitRankingRequest {
  rankings: RankingRequest[];
}

export interface TierResult {
  item_id: string;
  item_name: string;
  image_url: string;
  counts: Record<string, number>;
  total: number;
  top_tier: string;
  average_score: number;
  rank: number;
}

export interface TierlistResultResponse {
  tierlist: Tierlist;
  total_submissions: number;
  results: TierResult[];
}
