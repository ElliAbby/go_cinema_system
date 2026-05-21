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
