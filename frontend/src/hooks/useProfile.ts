import { useQuery } from "@tanstack/react-query";
import { authApi } from "../api/auth";
import { useAuth } from "../context/AuthContext";

export const useProfile = () => {
  const { user, isAuthenticated } = useAuth();

  return useQuery({
    queryKey: ["profile", user?.id],
    queryFn: () => authApi.getProfile(user!.id),
    enabled: isAuthenticated && !!user?.id,
    staleTime: 2 * 60 * 1000,
  });
};
