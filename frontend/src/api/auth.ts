import client from "./client";
import { AuthResponse, User } from "../types";

export const authApi = {
  register: async (
    email: string,
    password: string,
    phone?: string,
  ): Promise<AuthResponse> => {
    const response = await client.post("/auth/register", {
      email,
      password,
      phone,
    });
    return response.data;
  },

  login: async (email: string, password: string): Promise<AuthResponse> => {
    const response = await client.post("/auth/login", { email, password });
    return response.data;
  },

  getProfile: async (userID: number): Promise<User> => {
    const response = await client.get(`/users/${userID}`);
    return response.data;
  },
};
