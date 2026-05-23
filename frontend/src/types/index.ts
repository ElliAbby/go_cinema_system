export interface User {
  id: number;
  email: string;
  phone?: string;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface AuthResponse {
  token: string;
  user_id: number;
  email: string;
  expires_at: number;
}

export interface Movie {
  id: number;
  title: string;
  duration: number;
  rating: string;
  description: string;
  created_at?: string;
}

export interface Cinema {
  id: number;
  name: string;
  address: string;
  created_at?: string;
}

export interface Hall {
  id: number;
  name: string;
  cinema_id: number;
  rows: number;
  seats_per_row: number;
  total_seats: number;
}

export interface Session {
  id: number;
  movie_id: number;
  hall_id: number;
  start_time: string;
  price_base: number;
  available_seats?: number;
  created_at?: string;
}

export interface Booking {
  id: string;
  user_id: number;
  session_id?: number;
  total_price: number;
  status: "pending" | "paid" | "cancelled" | "refunded";
  seat_ids?: number[];
  created_at?: string;
}

export interface Ticket {
  id: string;
  session_id: number;
  session_start_time?: string;
  seat_id: number;
  booking_id: string;
  seat_number?: number;
  status: string;
  created_at?: string;
}

export interface CreateBookingPayload {
  session_id: number;
  seat_ids: number[];
}
