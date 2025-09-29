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
