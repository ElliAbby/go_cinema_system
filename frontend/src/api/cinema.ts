import client from "./client";
import { Movie, Cinema, Session, Hall } from "../types";

export const cinemaApi = {
  // Movies
  getMovies: async (): Promise<Movie[]> => {
    const response = await client.get("/movies");
    return response.data || [];
  },

  getMovie: async (id: number): Promise<Movie> => {
    const response = await client.get(`/movies/${id}`);
    return response.data;
  },

  // Cinemas
  getCinemas: async (): Promise<Cinema[]> => {
    const response = await client.get("/cinemas");
    return response.data || [];
  },

  getCinema: async (id: number): Promise<Cinema> => {
    const response = await client.get(`/cinemas/${id}`);
    return response.data;
  },

  getHalls: async (cinemaId: number): Promise<Hall[]> => {
    const response = await client.get(`/cinemas/${cinemaId}/halls`);
    return response.data || [];
  },

  getHall: async (id: number): Promise<Hall> => {
    const response = await client.get(`/halls/${id}`);
    return response.data;
  },

  // Sessions
  getSessions: async (): Promise<Session[]> => {
    const response = await client.get("/sessions");
    return response.data || [];
  },

  getSession: async (id: number): Promise<Session> => {
    const response = await client.get(`/sessions/${id}`);
    return response.data;
  },

  getMovieSessions: async (movieId: number): Promise<Session[]> => {
    const response = await client.get(`/movies/${movieId}/sessions`);
    return response.data || [];
  },

  // Seats
  getSeatsByHall: async (hallId: number): Promise<any[]> => {
    const response = await client.get(`/halls/${hallId}/seats`);
    return response.data || [];
  },

  getReservedSeats: async (sessionId: number): Promise<number[]> => {
    const response = await client.get(`/sessions/${sessionId}/reserved`);
    return response.data?.reserved_seat_ids || [];
  },
};
