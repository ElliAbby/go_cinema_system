import { useQuery } from "@tanstack/react-query";
import { cinemaApi } from "../api/cinema";

export const useSeatsByHall = (hallId: number | null) => {
  return useQuery({
    queryKey: ["seats", hallId],
    queryFn: () => cinemaApi.getSeatsByHall(hallId!),
    enabled: !!hallId,
    staleTime: 1 * 60 * 1000,
  });
};

export const useReservedSeats = (sessionId: number | null) => {
  return useQuery({
    queryKey: ["reservedSeats", sessionId],
    queryFn: () => cinemaApi.getReservedSeats(sessionId!),
    enabled: !!sessionId,
    staleTime: 10 * 1000,
  });
};
