import { useQuery } from "@tanstack/react-query";
import { cinemaApi } from "../api/cinema";

export const useCinemas = () => {
  return useQuery({
    queryKey: ["cinemas"],
    queryFn: () => cinemaApi.getCinemas(),
    staleTime: 5 * 60 * 1000,
  });
};

export const useCinema = (id: number | null) => {
  return useQuery({
    queryKey: ["cinema", id],
    queryFn: () => cinemaApi.getCinema(id!),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  });
};

export const useHalls = (cinemaId: number | null) => {
  return useQuery({
    queryKey: ["halls", cinemaId],
    queryFn: () => cinemaApi.getHalls(cinemaId!),
    enabled: !!cinemaId,
    staleTime: 5 * 60 * 1000,
  });
};

export const useHall = (id: number | null) => {
  return useQuery({
    queryKey: ["hall", id],
    queryFn: () => cinemaApi.getHall(id!),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  });
};

export const useHallSessions = (hallId: number | null) => {
  return useQuery({
    queryKey: ["hallSessions", hallId],
    queryFn: () => cinemaApi.getHallSessions(hallId!),
    enabled: !!hallId,
    staleTime: 1 * 60 * 1000,
  });
};
