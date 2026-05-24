import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  useBooking,
  useCancelBooking,
  usePurchaseBooking,
} from "../hooks/useBookings";
import { useCinema, useHall } from "../hooks/useCinemas";
import { useMovie } from "../hooks/useMovies";
import { useSession } from "../hooks/useSessions";
import { useSeatsByHall } from "../hooks/useSeats";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const BookingConfirm: React.FC = () => {
  const { bookingId } = useParams<{ bookingId: string }>();
  const navigate = useNavigate();
  const [isPaying, setIsPaying] = useState(false);

  const id = bookingId || null;
  const { data: booking, isLoading } = useBooking(id);
  const sessionId = booking?.session_id ?? null;
  const { data: session, isLoading: isSessionLoading } = useSession(sessionId);
  const { data: movie, isLoading: isMovieLoading } = useMovie(
    session?.movie_id ?? null,
  );
  const { data: hall, isLoading: isHallLoading } = useHall(
    session?.hall_id ?? null,
  );
  const { data: seats } = useSeatsByHall(hall?.id ?? null);
  const { data: cinema, isLoading: isCinemaLoading } = useCinema(
    hall?.cinema_id ?? null,
  );
  const { mutate: purchaseBooking, isPending: isProcessing } =
    usePurchaseBooking();
  const { mutate: cancelBooking, isPending: isCancelling } = useCancelBooking();

  if (
    isLoading ||
    isSessionLoading ||
    isMovieLoading ||
    isHallLoading ||
    isCinemaLoading
  ) {
    return <FullPageLoader />;
  }

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

  const handleCancel = () => {
    const cancellationMessage =
      booking.status === "paid"
        ? "Отменить это оплаченное бронирование? Билеты будут деактивированы, а места освободятся."
        : "Отменить это бронирование? Несохраненные места будут освобождены.";

    if (!window.confirm(cancellationMessage)) {
      return;
    }

    cancelBooking(booking.id, {
      onSuccess: () => {
        navigate("/bookings");
      },
    });
  };

  const statusColorMap = {
    pending: "text-yellow-300",
    paid: "text-emerald-300",
    cancelled: "text-red-300",
  };

  const formatDate = (value?: string) => {
    if (!value) return "—";

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;

    return new Intl.DateTimeFormat("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
    }).format(date);
  };

  const formatTime = (value?: string) => {
    if (!value) return "—";

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;

    return new Intl.DateTimeFormat("ru-RU", {
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  };

  const seatLabelById: Record<number, string> = {};
  (seats || []).forEach((seat) => {
    seatLabelById[seat.id] =
      `Ряд ${seat.row_number}, место ${seat.seat_number}`;
  });
  const seatNumbers = (booking.seat_ids || []).slice().sort((a, b) => a - b);
  const bookingDate = formatDate(booking.created_at);
  const sessionDate = formatDate(session?.start_time);
  const sessionTime = formatTime(session?.start_time);

  return (
    <div className="min-h-screen bg-dark-950 relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.22),_transparent_40%),radial-gradient(circle_at_bottom_right,_rgba(16,185,129,0.16),_transparent_35%)]" />

      <div className="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10 sm:py-14">
        <div className="overflow-hidden rounded-3xl border border-white/10 bg-white/5 shadow-2xl shadow-black/30 backdrop-blur-xl">
          <div className="bg-gradient-to-r from-blue-600 via-cyan-500 to-emerald-500 px-6 py-6 sm:px-8 sm:py-8">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
              <div>
                <div className="inline-flex items-center gap-2 rounded-full bg-black/20 px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em] text-white/90">
                  Бронирование
                </div>
                <h1 className="mt-4 text-3xl sm:text-4xl font-bold text-white">
                  Бронирование готово к оплате
                </h1>
                <p className="mt-2 max-w-2xl text-sm sm:text-base text-cyan-50/90">
                  Проверьте фильм, кинотеатр, зал и время сеанса перед оплатой.
                </p>
              </div>

              <div
                className={`inline-flex items-center gap-2 rounded-2xl border border-white/20 bg-black/20 px-4 py-3 text-sm font-semibold ${statusColorMap[booking.status as keyof typeof statusColorMap]}`}
              >
                {booking.status === "pending" && "⏳ Ожидание оплаты"}
                {booking.status === "paid" && "✓ Оплачено"}
                {booking.status === "cancelled" && "✗ Отменено"}
              </div>
            </div>
          </div>

          <div className="grid gap-6 px-6 py-6 sm:px-8 lg:grid-cols-[1.3fr_0.9fr]">
            <div className="space-y-6">
              <section className="rounded-2xl border border-dark-700 bg-dark-900/80 p-5 sm:p-6">
                <div className="flex flex-wrap items-start justify-between gap-4">
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
                      Фильм
                    </p>
                    <h2 className="mt-2 text-2xl font-bold text-white">
                      {movie?.title ?? "Фильм не найден"}
                    </h2>
                    <p className="mt-2 text-sm text-gray-400 line-clamp-2">
                      {movie?.description ?? "Информация о фильме недоступна"}
                    </p>
                  </div>

                  <div className="rounded-2xl border border-blue-500/20 bg-blue-500/10 px-4 py-3 text-right">
                    <p className="text-xs uppercase tracking-[0.2em] text-blue-200/70">
                      Начало
                    </p>
                    <p className="mt-1 text-lg font-semibold text-white">
                      {sessionTime}
                    </p>
                    <p className="text-sm text-blue-100/80">{sessionDate}</p>
                  </div>
                </div>

                <div className="mt-5 flex flex-wrap gap-2">
                  {movie?.rating && (
                    <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white">
                      {movie.rating}
                    </span>
                  )}
                  {movie?.duration && (
                    <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white">
                      {movie.duration} мин
                    </span>
                  )}
                  <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white">
                    {seatNumbers.length} мест
                  </span>
                </div>
              </section>

              <section className="grid gap-4 sm:grid-cols-3">
                <div className="rounded-2xl border border-dark-700 bg-dark-900 p-5">
                  <p className="text-sm text-gray-400">Кинотеатр</p>
                  <p className="mt-2 text-lg font-semibold text-white">
                    {cinema?.name ?? "Не найден"}
                  </p>
                  <p className="mt-1 text-sm text-gray-500">
                    {cinema?.address ?? "Адрес недоступен"}
                  </p>
                </div>

                <div className="rounded-2xl border border-dark-700 bg-dark-900 p-5">
                  <p className="text-sm text-gray-400">Зал</p>
                  <p className="mt-2 text-lg font-semibold text-white">
                    {hall?.name ?? "Не найден"}
                  </p>
                  <p className="mt-1 text-sm text-gray-500">
                    {hall?.total_seats ? `${hall.total_seats} мест` : ""}
                  </p>
                </div>

                <div className="rounded-2xl border border-dark-700 bg-dark-900 p-5">
                  <p className="text-sm text-gray-400">Время</p>
                  <p className="mt-2 text-lg font-semibold text-white">
                    {sessionTime}
                  </p>
                  <p className="mt-1 text-sm text-gray-500">{sessionDate}</p>
                </div>
              </section>
            </div>

            <aside className="space-y-4 rounded-2xl border border-dark-700 bg-dark-900/80 p-5 sm:p-6">
              <div className="rounded-2xl bg-dark-950/80 p-4">
                <p className="text-sm text-gray-400">Стоимость</p>
                <p className="mt-2 text-3xl font-bold text-emerald-400">
                  {Math.round(booking.total_price)} ₽
                </p>
                <p className="mt-2 text-sm text-gray-500">
                  Бронирование создано {bookingDate}
                </p>
              </div>

              <div className="rounded-2xl bg-dark-950/80 p-4">
                <p className="text-sm text-gray-400">Места</p>
                <div className="mt-3 flex flex-wrap gap-2">
                  {seatNumbers.map((seat) => (
                    <span
                      key={seat}
                      className="rounded-lg border border-blue-400/30 bg-blue-500/15 px-3 py-1 text-sm font-medium text-blue-100"
                    >
                      {seatLabelById[seat] ?? `Место ${seat}`}
                    </span>
                  ))}
                </div>
              </div>

              <div className="rounded-2xl bg-dark-950/80 p-4">
                <p className="text-sm text-gray-400">Статус брони</p>
                <p
                  className={`mt-2 text-lg font-semibold ${statusColorMap[booking.status as keyof typeof statusColorMap]}`}
                >
                  {booking.status === "pending" && "Ожидание оплаты"}
                  {booking.status === "paid" && "Оплачено"}
                  {booking.status === "cancelled" && "Отменено"}
                </p>
              </div>

              {(booking.status === "pending" || booking.status === "paid") && (
                <div className="space-y-3 pt-2">
                  {booking.status === "pending" && (
                    <button
                      onClick={handlePayment}
                      disabled={isProcessing}
                      className="w-full rounded-xl bg-emerald-500 px-4 py-3 font-semibold text-white transition hover:bg-emerald-400 disabled:cursor-not-allowed disabled:bg-emerald-500/50"
                    >
                      {isProcessing || isPaying
                        ? "Обработка платежа..."
                        : "Оплатить"}
                    </button>
                  )}

                  <button
                    onClick={handleCancel}
                    disabled={isCancelling || isProcessing || isPaying}
                    className="w-full rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 font-semibold text-red-100 transition hover:bg-red-500/20 disabled:cursor-not-allowed disabled:bg-red-500/5"
                  >
                    {isCancelling
                      ? "Отмена..."
                      : booking.status === "paid"
                        ? "Отменить и деактивировать билеты"
                        : "Отменить бронирование"}
                  </button>
                </div>
              )}

              <button
                onClick={() => navigate("/bookings")}
                className="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 font-semibold text-gray-200 transition hover:bg-white/10"
              >
                Мои бронирования
              </button>
            </aside>
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
