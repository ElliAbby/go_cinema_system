import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useSession } from "../hooks/useSessions";
import { useCreateBooking } from "../hooks/useBookings";
import { useSeatsByHall, useReservedSeats } from "../hooks/useSeats";
import { useAuth } from "../context/AuthContext";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { SeatSelector } from "../components/SeatSelector";

export const Booking: React.FC = () => {
  const { sessionId } = useParams<{ sessionId: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  // selectedSeats will store DB seat IDs
  const [selectedSeats, setSelectedSeats] = useState<Set<number>>(new Set());

  const id = sessionId ? parseInt(sessionId) : null;
  const { data: session, isLoading } = useSession(id);
  const {
    data: seats,
    isLoading: seatsLoading,
    error: seatsError,
  } = useSeatsByHall(session ? session.hall_id : null);
  const {
    data: reservedSeatIds,
    isLoading: reservedLoading,
    error: reservedError,
  } = useReservedSeats(session ? session.id : null);

  // Log for debugging
  React.useEffect(() => {
    if (session) {
      console.log("Session loaded:", session);
    }
  }, [session]);

  React.useEffect(() => {
    if (seats) {
      console.log("Seats loaded:", seats);
    }
  }, [seats]);

  React.useEffect(() => {
    if (reservedSeatIds) {
      console.log("Reserved seats:", reservedSeatIds);
    }
  }, [reservedSeatIds]);

  // build mapping from visual seatNumber -> DB seat id
  const seatIdMap: Record<number, number> = {};
  if (seats && seats.length > 0) {
    const maxSeatPerRow = Math.max(
      ...seats.map((s: any) => s.seat_number || 0),
    );
    seats.forEach((s: any) => {
      const seatNumber = (s.row_number - 1) * maxSeatPerRow + s.seat_number;
      seatIdMap[seatNumber] = s.id;
    });
  }
  const { mutate: createBooking, isPending: isCreating } = useCreateBooking();

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">
            Пожалуйста, войдите в аккаунт для бронирования
          </p>
          <button
            onClick={() => navigate("/login")}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Войти
          </button>
        </div>
      </div>
    );
  }

  if (isLoading || seatsLoading || reservedLoading) return <FullPageLoader />;

  if (!session) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Сеанс не найден</p>
          <button
            onClick={() => navigate("/")}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Вернуться на главную
          </button>
        </div>
      </div>
    );
  }

  if (seatsError || reservedError) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">
            Ошибка при загрузке мест:{" "}
            {seatsError?.message || reservedError?.message}
          </p>
          <button
            onClick={() => navigate("/")}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Вернуться на главную
          </button>
        </div>
      </div>
    );
  }

  if (!session) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Сеанс не найден</p>
          <button
            onClick={() => navigate("/")}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Вернуться на главную
          </button>
        </div>
      </div>
    );
  }

  const handleSeatClick = (_seatNumber: number, seatId?: number) => {
    if (!seatId) return;
    const newSelected = new Set(selectedSeats);
    if (newSelected.has(seatId)) {
      newSelected.delete(seatId);
    } else {
      newSelected.add(seatId);
    }
    setSelectedSeats(newSelected);
  };

  const handleBooking = async () => {
    if (selectedSeats.size === 0) return;

    createBooking(
      {
        session_id: session.id,
        seat_ids: Array.from(selectedSeats),
      },
      {
        onSuccess: (booking: any) => {
          console.log("Booking created:", booking);
          navigate(`/bookings/${booking.id}/purchase`);
        },
        onError: (error: any) => {
          console.error("Booking error:", error);
          alert(
            "Ошибка при бронировании: " +
              (error?.response?.data?.error?.message ||
                error.message ||
                "Неизвестная ошибка"),
          );
        },
      },
    );
  };

  const totalPrice = selectedSeats.size * session.price_base;

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <button
          onClick={() => navigate(-1)}
          className="mb-8 text-blue-400 hover:text-blue-300 transition"
        >
          ← Назад
        </button>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Seats */}
          <div className="lg:col-span-2">
            <h1 className="text-3xl font-bold text-white mb-8">
              Выберите места
            </h1>
            <SeatSelector
              rows={10}
              seatsPerRow={10}
              selectedSeats={selectedSeats}
              onSeatClick={handleSeatClick}
              seatIdMap={seatIdMap}
              reservedSeatIds={new Set(reservedSeatIds || [])}
            />
          </div>

          {/* Summary */}
          <div className="lg:col-span-1">
            <div className="sticky top-24 bg-dark-800 border border-dark-700 rounded-lg p-6 space-y-6">
              <div>
                <h3 className="text-lg font-semibold text-white mb-2">Итого</h3>
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between text-gray-400">
                  <span>Цена за билет:</span>
                  <span>{Math.round(session.price_base)} ₽</span>
                </div>
                <div className="flex justify-between text-gray-400">
                  <span>Количество мест:</span>
                  <span>{selectedSeats.size}</span>
                </div>
              </div>

              <div className="border-t border-dark-700 pt-4">
                <div className="flex justify-between items-baseline">
                  <span className="text-gray-300">Итого:</span>
                  <span className="text-3xl font-bold text-blue-400">
                    {Math.round(totalPrice)} ₽
                  </span>
                </div>
              </div>

              {selectedSeats.size > 0 && (
                <div>
                  <p className="text-xs text-gray-400 mb-2">Выбранные места:</p>
                  <div className="flex flex-wrap gap-2">
                    {(() => {
                      const seatNumberById: Record<number, number> = {};
                      Object.entries(seatIdMap).forEach(([numStr, id]) => {
                        seatNumberById[id] = parseInt(numStr, 10);
                      });
                      return Array.from(selectedSeats)
                        .sort((a, b) => a - b)
                        .map((seatId) => (
                          <span
                            key={seatId}
                            className="px-2 py-1 bg-blue-500 text-white rounded text-xs"
                          >
                            {seatNumberById[seatId] ?? seatId}
                          </span>
                        ));
                    })()}
                  </div>
                </div>
              )}

              <button
                onClick={handleBooking}
                disabled={selectedSeats.size === 0 || isCreating}
                className="w-full py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-semibold rounded-lg transition"
              >
                {isCreating ? "Загрузка..." : "Продолжить к оплате"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const FullPageLoader = () => (
  <div className="min-h-screen bg-dark-950 flex items-center justify-center">
    <LoadingSpinner />
  </div>
);
