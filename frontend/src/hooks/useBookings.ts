import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { bookingApi } from "../api/booking";
import { CreateBookingPayload } from "../types";

export const useBookings = () => {
  return useQuery({
    queryKey: ["bookings"],
    queryFn: () => bookingApi.getBookings(),
    staleTime: 1 * 60 * 1000,
  });
};

export const useBooking = (id: string | null) => {
  return useQuery({
    queryKey: ["booking", id],
    queryFn: () => bookingApi.getBooking(id!),
    enabled: !!id,
    staleTime: 30 * 1000,
  });
};

export const useCreateBooking = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateBookingPayload) => bookingApi.createBooking(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bookings"] });
    },
  });
};

export const usePurchaseBooking = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => bookingApi.purchaseBooking(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bookings"] });
    },
  });
};

export const useTickets = () => {
  return useQuery({
    queryKey: ["tickets"],
    queryFn: () => bookingApi.getTickets(),
    staleTime: 5 * 60 * 1000,
  });
};
