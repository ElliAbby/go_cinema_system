import { useQuery } from "@tanstack/react-query";
import { cinemaApi } from "../api/cinema";

export const useSessions = () => {
  return useQuery({
    queryKey: ["sessions"],
    queryFn: () => cinemaApi.getSessions(),
    staleTime: 1 * 60 * 1000,
  });
};

export const useSession = (id: number | null) => {
  return useQuery({
    queryKey: ["session", id],
    queryFn: () => cinemaApi.getSession(id!),
    enabled: !!id,
    staleTime: 1 * 60 * 1000,
  });
};
