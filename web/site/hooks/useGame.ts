"use client";

import { useQuery, useMutation, UseQueryOptions, UseMutationOptions, useQueryClient } from "@tanstack/react-query";
import { api } from "../lib/api";
import type {
  GameEggDTO,
  CreateGameEggPayload,
  InventoryItemDTO,
  PlayerAccountDTO,
  PlayerEquipmentDTO,
} from "../lib/types";

// Keys
const qk = {
  eggs: (player_id: string) => ["game-eggs", player_id] as const,
  inventory: (player_id: string) => ["game-inventory", player_id] as const,
  playerByAccount: (account_id: string) => ["game-player-account", account_id] as const,
  playerEquipment: (player_id: string) => ["game-player-equipment", player_id] as const,
};

export function useGameEggs(player_id: string, options?: Partial<UseQueryOptions<GameEggDTO[]>>) {
  return useQuery({
    queryKey: qk.eggs(player_id),
    queryFn: () => api.getGameEggs(player_id),
    enabled: !!player_id,
    staleTime: 60_000,
    ...options,
  });
}

export function useCreateGameEgg(options?: UseMutationOptions<unknown, Error, CreateGameEggPayload>) {
  const qc = useQueryClient();
  return useMutation<unknown, Error, CreateGameEggPayload>({
    mutationFn: (payload) => api.createGameEgg(payload),
    onSuccess: (_data, variables) => {
      if (variables.player_id) qc.invalidateQueries({ queryKey: qk.eggs(variables.player_id) });
    },
    ...options,
  });
}

export function useInventory(player_id: string, options?: Partial<UseQueryOptions<InventoryItemDTO[]>>) {
  return useQuery({
    queryKey: qk.inventory(player_id),
    queryFn: () => api.getInventory(player_id),
    enabled: !!player_id,
    staleTime: 60_000,
    ...options,
  });
}

export function usePlayerByAccount(account_id: string, options?: Partial<UseQueryOptions<PlayerAccountDTO>>) {
  return useQuery({
    queryKey: qk.playerByAccount(account_id),
    queryFn: () => api.getPlayerByAccount(account_id),
    enabled: !!account_id,
    staleTime: 60_000,
    ...options,
  });
}

export function usePlayerEquipment(player_id: string, options?: Partial<UseQueryOptions<PlayerEquipmentDTO[]>>) {
  return useQuery({
    queryKey: qk.playerEquipment(player_id),
    queryFn: () => api.getPlayerEquipment(player_id),
    enabled: !!player_id,
    staleTime: 60_000,
    ...options,
  });
}
