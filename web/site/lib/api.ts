export type ApiOptions = {
  baseUrl?: string;
  headers?: Record<string, string>;
};

// Align with server base path
const DEFAULT_BASE_URL = "http://localhost:8080/api/v1";

async function request<T>(
  path: string,
  options: RequestInit = {},
  { baseUrl = DEFAULT_BASE_URL, headers = {} }: ApiOptions = {}
): Promise<T> {
  const url = `${baseUrl}${path}`;
  const res = await fetch(url, {
    credentials: 'include', // include cookies set by the server
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...headers,
      ...(options.headers as Record<string, string>),
    },
  });

  const contentType = res.headers.get('content-type') || '';
  const isJson = contentType.includes('application/json');
  const body = isJson ? await res.json().catch(() => ({})) : await res.text();

  if (!res.ok) {
    const message = isJson && body && (body.message || body.error) ? (body.message || body.error) : res.statusText;
    throw new Error(message || `Request failed: ${res.status}`);
  }

  return body as T;
}

export type RegisterPayload = {
  email: string;
  password: string;
  username?: string;
};

export type LoginPayload = {
  email: string;
  password: string;
};

// Aligning to AccountLoginResponse subset the frontend needs
export type AuthResponse = {
  access_token?: string;
  access_token_expires_at?: string;
  refresh_token?: string;
  refresh_token_expires_at?: string;
  user?: {
    id?: string;
    email?: string;
    first_name?: string;
    last_name?: string;
    user_name?: string;
    role?: string;
  };
};

export const api = {
  post: request as <T>(path: string, options?: RequestInit, opts?: ApiOptions) => Promise<T>,
  register(payload: RegisterPayload, opts?: ApiOptions) {
    return request<AuthResponse>("/auth/register", { method: 'POST', body: JSON.stringify(payload) }, opts);
  },
  login(payload: LoginPayload, opts?: ApiOptions) {
    return request<AuthResponse>("/auth/login", { method: 'POST', body: JSON.stringify(payload) }, opts);
  },
  renew(opts?: ApiOptions) {
    return request<AuthResponse>("/auth/renew", { method: 'POST' }, opts);
  },
  // ---- Game endpoints ----
  getGameEggs(player_id: string, opts?: ApiOptions) {
    const qp = new URLSearchParams({ player_id }).toString();
    return request<import('./types').GameEggDTO[]>(`/game/eggs?${qp}`, { method: 'GET' }, opts);
  },
  createGameEgg(payload: import('./types').CreateGameEggPayload, opts?: ApiOptions) {
    return request<unknown>(`/game/eggs`, { method: 'POST', body: JSON.stringify(payload) }, opts);
  },
  getInventory(player_id: string, opts?: ApiOptions) {
    const qp = new URLSearchParams({ player_id }).toString();
    return request<import('./types').InventoryItemDTO[]>(`/game/inventory?${qp}`, { method: 'GET' }, opts);
  },
  getPlayerByAccount(account_id: string, opts?: ApiOptions) {
    const qp = new URLSearchParams({ account_id }).toString();
    return request<import('./types').PlayerAccountDTO>(`/game/player?${qp}`, { method: 'GET' }, opts);
  },
  getPlayerEquipment(player_id: string, opts?: ApiOptions) {
    const qp = new URLSearchParams({ player_id }).toString();
    return request<import('./types').PlayerEquipmentDTO[]>(`/game/player?${qp}`, { method: 'GET' }, opts);
  },
};
