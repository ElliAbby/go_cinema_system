import { useQuery } from "@tanstack/react-query";
import { cinemaApi } from "../api/cinema";

export const useMovies = () => {
  return useQuery({
    queryKey: ["movies"],
    queryFn: () => cinemaApi.getMovies(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useMovie = (id: number | null) => {
  return useQuery({
    queryKey: ["movie", id],
    queryFn: () => cinemaApi.getMovie(id!),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  });
};

export const useMovieSessions = (movieId: number | null) => {
  return useQuery({
    queryKey: ["movieSessions", movieId],
    queryFn: () => cinemaApi.getMovieSessions(movieId!),
    enabled: !!movieId,
    staleTime: 1 * 60 * 1000,
  });
};
