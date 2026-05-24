import React from "react";
import { useNavigate } from "react-router-dom";
import { useBookings, useCancelBooking } from "../hooks/useBookings";
import { useAuth } from "../context/AuthContext";
import { LoadingSpinner } from "../components/LoadingSpinner";

type BookingStatus = "pending" | "paid" | "cancelled" | "refunded";

const statusMeta: Record<
  BookingStatus,
  { label: string; badgeClass: string; dotClass: string }
> = {
  pending: {
    label: "Ожидание",
    badgeClass: "bg-amber-500/10 text-amber-300 border-amber-500/20",
    dotClass: "bg-amber-400",
  },
  paid: {
    label: "Оплачено",
    badgeClass: "bg-emerald-500/10 text-emerald-300 border-emerald-500/20",
    dotClass: "bg-emerald-400",
  },
  cancelled: {
    label: "Отменено",
    badgeClass: "bg-rose-500/10 text-rose-300 border-rose-500/20",
    dotClass: "bg-rose-400",
  },
  refunded: {
    label: "Возврат",
    badgeClass: "bg-slate-500/10 text-slate-300 border-slate-500/20",
    dotClass: "bg-slate-400",
  },
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

const BookingCard: React.FC<{
  booking: {
    id: string;
    created_at?: string;
    total_price: number;
    status: BookingStatus;
  };
  isCancelling: boolean;
  onOpen: () => void;
  onCancel: (event: React.MouseEvent) => void;
}> = ({ booking, isCancelling, onOpen, onCancel }) => {
  const meta = statusMeta[booking.status] ?? statusMeta.cancelled;

  return (
    <article
      onClick={onOpen}
      className="group cursor-pointer rounded-3xl border border-white/8 bg-white/[0.03] p-5 shadow-[0_18px_50px_rgba(0,0,0,0.18)] transition hover:border-white/12 hover:bg-white/[0.05]"
    >
      <div className="flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <span
              className={`inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium ${meta.badgeClass}`}
            >
              <span className={`h-2 w-2 rounded-full ${meta.dotClass}`} />
              {meta.label}
            </span>
            <span className="rounded-full border border-white/8 bg-black/20 px-3 py-1 text-xs text-gray-300">
              Заказ #{booking.id.slice(0, 8)}
            </span>
          </div>

          <div className="mt-4 space-y-1">
            <p className="text-sm uppercase tracking-[0.2em] text-gray-500">
              Дата заказа
            </p>
            <p className="text-lg font-semibold text-white">
              {formatDate(booking.created_at)}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-4 sm:justify-end">
          <div className="text-right">
            <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
              Сумма
            </p>
            <p className="mt-1 text-3xl font-semibold text-white">
              {Math.round(booking.total_price)} ₽
            </p>
          </div>

          <div className="flex flex-col items-stretch gap-2">
            <button
              onClick={onCancel}
              disabled={isCancelling}
              className="rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm font-medium text-gray-200 transition hover:bg-white/10 disabled:cursor-not-allowed disabled:bg-white/5 disabled:text-gray-500"
            >
              Отменить
            </button>
            <span className="text-center text-xs text-gray-500 opacity-0 transition group-hover:opacity-100">
              Открыть заказ
            </span>
          </div>
        </div>
      </div>
    </article>
  );
};

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
    <div className="min-h-screen bg-dark-950 relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.16),_transparent_40%),radial-gradient(circle_at_bottom_right,_rgba(16,185,129,0.1),_transparent_35%)]" />

      <div className="relative mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
        <div className="mb-10 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.24em] text-gray-500">
              Orders
            </p>
            <h1 className="mt-2 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
              Мои заказы
            </h1>
            <p className="mt-2 max-w-2xl text-sm text-gray-400 sm:text-base">
              Все бронирования, их статусы и быстрый доступ к оплате или отмене.
            </p>
          </div>
        </div>

        {isLoading ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <LoadingSpinner />
          </div>
        ) : !bookings || bookings.length === 0 ? (
          <div className="rounded-3xl border border-white/8 bg-white/[0.03] px-6 py-16 text-center">
            <p className="text-lg font-medium text-white">
              У вас пока нет заказов
            </p>
            <p className="mt-2 text-sm text-gray-400">
              Когда вы оформите бронирование, оно появится здесь.
            </p>
            <button
              onClick={() => navigate("/")}
              className="mt-6 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-blue-500"
            >
              Выбрать фильм
            </button>
          </div>
        ) : (
          <div className="grid gap-4">
            {bookings.map((booking) => {
              const canCancel =
                booking.status === "pending" || booking.status === "paid";

              return (
                <BookingCard
                  key={booking.id}
                  booking={booking as any}
                  isCancelling={isCancelling}
                  onOpen={() => navigate(`/bookings/${booking.id}/purchase`)}
                  onCancel={(event) => {
                    event.stopPropagation();

                    if (!canCancel) return;
                    handleCancel(booking.id, event);
                  }}
                />
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};
