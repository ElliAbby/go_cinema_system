import client from "./client";
import { Booking, Ticket, CreateBookingPayload } from "../types";

export const bookingApi = {
  getBookings: async (): Promise<Booking[]> => {
    const response = await client.get("/bookings");
    return response.data || [];
  },

  getBooking: async (id: string): Promise<Booking> => {
    const response = await client.get(`/bookings/${id}`);
    return response.data;
  },

  createBooking: async (payload: CreateBookingPayload): Promise<Booking> => {
    const response = await client.post("/bookings", payload);
    return response.data.booking || response.data;
  },

  purchaseBooking: async (
    id: string,
  ): Promise<{ booking?: Booking; message?: string }> => {
    const response = await client.post(`/bookings/${id}/purchase`);
    return response.data;
  },

  cancelBooking: async (
    id: string,
  ): Promise<{ booking?: Booking; message?: string }> => {
    const response = await client.post(`/bookings/${id}/cancel`);
    return response.data;
  },

  getTickets: async (): Promise<Ticket[]> => {
    const response = await client.get("/tickets");
    return response.data || [];
  },

  getTicket: async (id: number): Promise<Ticket> => {
    const response = await client.get(`/tickets/${id}`);
    return response.data;
  },
};
