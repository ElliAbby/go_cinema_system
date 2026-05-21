import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useBooking, usePurchaseBooking } from "../hooks/useBookings";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const BookingConfirm: React.FC = () => {
  const { bookingId } = useParams<{ bookingId: string }>();
  const navigate = useNavigate();
  const [isPaying, setIsPaying] = useState(false);

  const id = bookingId || null;
  const { data: booking, isLoading } = useBooking(id);
  const { mutate: purchaseBooking, isPending: isProcessing } =
    usePurchaseBooking();

  if (isLoading) return <FullPageLoader />;

  if (!booking) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Бронирование не найдено</p>
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

  const handlePayment = async () => {
    setIsPaying(true);
    purchaseBooking(booking.id, {
      onSuccess: () => {
        setIsPaying(false);
        navigate("/bookings");
      },
      onError: () => {
        setIsPaying(false);
      },
    });
  };

  const statusColorMap = {
    pending: "text-yellow-400",
    paid: "text-green-400",
    cancelled: "text-red-400",
  };

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-dark-800 border border-dark-700 rounded-lg p-8">
          <div className="text-center mb-8">
            <div className="text-4xl mb-4">✅</div>
            <h1 className="text-3xl font-bold text-white mb-2">
              Бронирование создано
            </h1>
            <p className="text-gray-400">Номер бронирования: #{booking.id}</p>
          </div>

          <div className="space-y-6 mb-8">
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-dark-900 rounded-lg p-4">
                <p className="text-gray-400 text-sm mb-1">Количество мест</p>
                <p className="text-2xl font-bold text-blue-400">
                  {(booking.seat_ids || []).length}
                </p>
              </div>
              <div className="bg-dark-900 rounded-lg p-4">
                <p className="text-gray-400 text-sm mb-1">Стоимость</p>
                <p className="text-2xl font-bold text-green-400">
                  {Math.round(booking.total_price)} ₽
                </p>
              </div>
            </div>

            <div className="bg-dark-900 rounded-lg p-4">
              <p className="text-gray-400 text-sm mb-2">Статус</p>
              <div
                className={`text-lg font-semibold ${statusColorMap[booking.status as keyof typeof statusColorMap]}`}
              >
                {booking.status === "pending" && "⏳ Ожидание оплаты"}
                {booking.status === "paid" && "✓ Оплачено"}
                {booking.status === "cancelled" && "✗ Отменено"}
              </div>
            </div>

            <div className="bg-dark-900 rounded-lg p-4">
              <p className="text-gray-400 text-sm mb-2">Выбранные места</p>
              <div className="flex flex-wrap gap-2">
                {(booking.seat_ids || [])
                  .slice()
                  .sort((a, b) => a - b)
                  .map((seat) => (
                    <span
                      key={seat}
                      className="px-3 py-1 bg-blue-500 text-white rounded-lg text-sm font-medium"
                    >
                      {seat}
                    </span>
                  ))}
              </div>
            </div>
          </div>

          {booking.status === "pending" && (
            <button
              onClick={handlePayment}
              disabled={isProcessing}
              className="w-full py-3 bg-green-600 hover:bg-green-700 disabled:bg-green-600/50 text-white font-semibold rounded-lg transition mb-4"
            >
              {isProcessing || isPaying ? "Обработка платежа..." : "Оплатить"}
            </button>
          )}

          <button
            onClick={() => navigate("/bookings")}
            className="w-full py-3 bg-dark-700 hover:bg-dark-600 text-gray-300 font-semibold rounded-lg transition"
          >
            Мои бронирования
          </button>
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
