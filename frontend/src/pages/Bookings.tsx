import React from "react";
import { useNavigate } from "react-router-dom";
import { useBookings, useCancelBooking } from "../hooks/useBookings";
import { useAuth } from "../context/AuthContext";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const Bookings: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const { data: bookings, isLoading } = useBookings();
  const { mutate: cancelBooking, isPending: isCancelling } = useCancelBooking();

  const handleCancel = (bookingId: string, event: React.MouseEvent) => {
    event.stopPropagation();

    if (!window.confirm("Отменить это бронирование?")) {
      return;
    }

    cancelBooking(bookingId);
  };

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Пожалуйста, войдите в аккаунт</p>
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

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-4xl font-bold text-white mb-2">Мои бронирования</h1>
        <p className="text-gray-400 mb-12">История всех ваших бронирований</p>

        {isLoading ? (
          <LoadingSpinner />
        ) : !bookings || bookings.length === 0 ? (
          <div className="text-center py-12 bg-dark-800 rounded-lg border border-dark-700">
            <p className="text-gray-400 mb-4">У вас нет бронирований</p>
            <button
              onClick={() => navigate("/")}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
            >
              Выбрать фильм
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {bookings.map((booking) => {
              const canCancel =
                booking.status === "pending" || booking.status === "paid";

              return (
                <div
                  key={booking.id}
                  onClick={() => navigate(`/bookings/${booking.id}/purchase`)}
                  className="bg-dark-800 border border-dark-700 rounded-lg p-6 hover:border-dark-600 transition cursor-pointer"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <p className="text-lg font-semibold text-white mb-2">
                        Бронирование
                      </p>
                      <div className="space-y-1 text-sm text-gray-400">
                        <p>
                          Дата:{" "}
                          {new Date(
                            booking.created_at || "",
                          ).toLocaleDateString("ru-RU")}
                        </p>
                      </div>
                    </div>

                    <div className="text-right space-y-2">
                      <p className="text-2xl font-bold text-blue-400">
                        {Math.round(booking.total_price)} ₽
                      </p>
                      <div
                        className={`text-sm font-medium px-3 py-1 rounded-lg ${
                          booking.status === "pending"
                            ? "bg-yellow-900/30 text-yellow-400"
                            : booking.status === "paid"
                              ? "bg-green-900/30 text-green-400"
                              : "bg-red-900/30 text-red-400"
                        }`}
                      >
                        {booking.status === "pending" && "⏳ Ожидание"}
                        {booking.status === "paid" && "✓ Оплачено"}
                        {booking.status === "cancelled" && "✗ Отменено"}
                      </div>
                      {canCancel && (
                        <button
                          onClick={(event) => handleCancel(booking.id, event)}
                          disabled={isCancelling}
                          className="mt-2 px-3 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-600/50 text-white text-sm font-medium rounded-lg transition"
                        >
                          Отменить
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};
