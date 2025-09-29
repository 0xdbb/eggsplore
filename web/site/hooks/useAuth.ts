"use client";

import { useMutation } from "@tanstack/react-query";
import { api, LoginPayload, RegisterPayload, AuthResponse } from "../lib/api";

export function useRegister() {
  return useMutation<AuthResponse, Error, RegisterPayload>({
    mutationFn: (payload: RegisterPayload) => api.register(payload),
  });
}

export function useLogin() {
  return useMutation<AuthResponse, Error, LoginPayload>({
    mutationFn: (payload: LoginPayload) => api.login(payload),
  });
}
