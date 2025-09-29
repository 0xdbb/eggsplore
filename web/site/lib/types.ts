export type ID = string;

export type Coordinates = {
  latitude: number;
  longitude: number;
};

export type User = {
  id: ID;
  email: string;
  username?: string;
  firstName?: string;
  lastName?: string;
  role?: string;
  xp: number;
  coins?: number;
  createdAt?: string;
  updatedAt?: string;
  lastActive?: string;
};

export type Egg = {
  id: ID;
  ownerId: ID;
  title?: string;
  description?: string;
  coords: Coordinates;
  createdAt: string;
  color?: string;
  rarity?: "common" | "rare" | "epic" | "legendary";
};

export type Session = {
  accessToken?: string;
  refreshToken?: string;
  accessTokenExpiresAt?: string;
  refreshTokenExpiresAt?: string;
};

// ---- Game API DTOs ----
export type GameEggDTO = {
  collected_at: string;
  hatched: boolean;
  inventory_id: string;
  message: string;
  type: string; // e.g. "BUNNY" | "DRAGON"
};

export type CreateGameEggPayload = {
  lat: number;
  lon: number;
  message: string;
  player_id: string;
  type: "BUNNY" | "DRAGON" | string;
};

export type InventoryItemDTO = {
  created_at: string;
  description: string;
  id: string;
  item_type: string;
  player_id: string;
  quantity: number;
};

export type PlayerAccountDTO = {
  account_id: string;
  coins: number;
  created_at: string;
  id: string;
  level: number;
  settings: number[];
  updated_at: string;
  xp: number;
};

export type PlayerEquipmentDTO = {
  description: string;
  durability: number;
  equipped: boolean;
  inventory_id: string;
};
