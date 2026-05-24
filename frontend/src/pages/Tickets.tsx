import React from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useTickets } from "../hooks/useBookings";
import { useCinema, useHall } from "../hooks/useCinemas";
import { useMovie } from "../hooks/useMovies";
import { useSession } from "../hooks/useSessions";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { Ticket as TicketType } from "../types";

type TicketStatus = "active" | "deactivated" | "used";

const statusMeta: Record<
  TicketStatus,
  { label: string; badgeClass: string; dotClass: string }
> = {
  active: {
    label: "Активный",
    badgeClass: "bg-emerald-500/10 text-emerald-300 border-emerald-500/20",
    dotClass: "bg-emerald-400",
  },
  deactivated: {
    label: "Деактивированный",
    badgeClass: "bg-amber-500/10 text-amber-300 border-amber-500/20",
    dotClass: "bg-amber-400",
  },
  used: {
    label: "Использованный",
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

const formatTime = (value?: string) => {
  if (!value) return "—";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("ru-RU", {
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
};

const TicketCard: React.FC<{ ticket: TicketType }> = ({ ticket }) => {
  const sessionId = ticket.session_id ?? null;
  const { data: session } = useSession(sessionId);
  const { data: movie } = useMovie(session?.movie_id ?? null);
  const { data: hall } = useHall(session?.hall_id ?? null);
  const { data: cinema } = useCinema(hall?.cinema_id ?? null);

  const meta = statusMeta[ticket.status as TicketStatus] || statusMeta.used;
  const startTime = ticket.session_start_time || session?.start_time;
  const seatLabel =
    ticket.row_number && ticket.seat_number
      ? `Ряд ${ticket.row_number}, место ${ticket.seat_number}`
      : `Место ${ticket.seat_id}`;

  return (
    <article className="group rounded-3xl border border-white/8 bg-white/[0.03] p-5 shadow-[0_18px_50px_rgba(0,0,0,0.22)] transition hover:border-white/12 hover:bg-white/[0.05]">
      <div className="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0 flex-1">
          <div className="flex flex-wrap items-center gap-2">
            <span
              className={`inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium ${meta.badgeClass}`}
            >
              <span className={`h-2 w-2 rounded-full ${meta.dotClass}`} />
              {meta.label}
            </span>
            <span className="rounded-full border border-white/8 bg-black/20 px-3 py-1 text-xs text-gray-300">
              {seatLabel}
            </span>
          </div>

          <h3 className="mt-4 truncate text-xl font-semibold text-white">
            {movie?.title ?? "Фильм недоступен"}
          </h3>
          <p className="mt-1 text-sm text-gray-400">
            {cinema?.name ?? "Кинотеатр недоступен"}
          </p>
          <p className="mt-1 text-sm text-gray-500">
            {hall?.name ? `${hall.name} · ` : ""}
            {startTime
              ? `${formatDate(startTime)} · ${formatTime(startTime)}`
              : "Время сеанса недоступно"}
          </p>
        </div>

        <div className="rounded-2xl border border-white/8 bg-black/20 px-4 py-3 text-right sm:min-w-[140px]">
          <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
            Билет
          </p>
          <p className="mt-2 text-2xl font-semibold text-white">
            {ticket.id.slice(0, 8)}
          </p>
        </div>
      </div>

      <div className="mt-5 grid gap-3 sm:grid-cols-3">
        <div className="rounded-2xl border border-white/8 bg-black/20 p-4">
          <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
            Сеанс
          </p>
          <p className="mt-2 text-sm font-medium text-white">
            {startTime ? formatTime(startTime) : "—"}
          </p>
          <p className="mt-1 text-xs text-gray-500">
            {startTime ? formatDate(startTime) : "—"}
          </p>
        </div>

        <div className="rounded-2xl border border-white/8 bg-black/20 p-4">
          <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
            Зал
          </p>
          <p className="mt-2 text-sm font-medium text-white">
            {hall?.name ?? "—"}
          </p>
          <p className="mt-1 text-xs text-gray-500">{cinema?.address ?? "—"}</p>
        </div>

        <div className="rounded-2xl border border-white/8 bg-black/20 p-4">
          <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
            Фильм
          </p>
          <p className="mt-2 text-sm font-medium text-white">
            {movie?.title ?? "—"}
          </p>
          <p className="mt-1 text-xs text-gray-500">
            {movie?.duration ? `${movie.duration} мин` : "—"}
          </p>
        </div>
      </div>
    </article>
  );
};

const Section: React.FC<{
  title: string;
  description: string;
  tickets: TicketType[];
  accent: string;
  defaultOpen?: boolean;
}> = ({ title, description, tickets, accent, defaultOpen }) => {
  if (tickets.length === 0) return null;

  return (
    <details
      className="group rounded-3xl border border-white/8 bg-white/[0.03] shadow-[0_18px_50px_rgba(0,0,0,0.18)] open:border-white/12"
      open={defaultOpen}
    >
      <summary className="cursor-pointer list-none px-5 py-4 sm:px-6">
        <div className="flex items-center justify-between gap-4">
          <div className="min-w-0">
            <div className="flex items-center gap-3">
              <span className={`h-2.5 w-2.5 rounded-full ${accent}`} />
              <h2 className="text-xl font-semibold text-white sm:text-2xl">
                {title}
              </h2>
            </div>
            <p className="mt-1 text-sm text-gray-400">{description}</p>
          </div>

          <div className="flex items-center gap-3">
            <span className="rounded-full border border-white/10 bg-black/20 px-3 py-1 text-xs font-medium text-gray-300">
              {tickets.length}
            </span>
            <span className="text-gray-500 transition group-open:rotate-180">
              ▾
            </span>
          </div>
        </div>
      </summary>

      <div className="px-5 pb-5 sm:px-6">
        <div className="max-h-[34rem] space-y-4 overflow-y-auto pr-1 scrollbar-thin scrollbar-thumb-white/10 scrollbar-track-transparent">
          {tickets.map((ticket) => (
            <TicketCard key={ticket.id} ticket={ticket} />
          ))}
        </div>
      </div>
    </details>
  );
};

export const Tickets: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const { data: tickets, isLoading, isError, error } = useTickets();

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

  const activeTickets = (tickets || []).filter(
    (ticket) => ticket.status === "active",
  );
  const deactivatedTickets = (tickets || []).filter(
    (ticket) => ticket.status === "deactivated",
  );
  const usedTickets = (tickets || []).filter(
    (ticket) => ticket.status === "used",
  );

  return (
    <div className="min-h-screen bg-dark-950 relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.18),_transparent_40%),radial-gradient(circle_at_bottom_right,_rgba(16,185,129,0.12),_transparent_35%)]" />

      <div className="relative mx-auto max-w-7xl px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
        <div className="mb-10 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.24em] text-gray-500">
              Tickets
            </p>
            <h1 className="mt-2 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
              Мои билеты
            </h1>
            <p className="mt-2 max-w-2xl text-sm text-gray-400 sm:text-base">
              Все ваши билеты в одном месте.
            </p>
          </div>
        </div>

        {isLoading ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <LoadingSpinner />
          </div>
        ) : isError ? (
          <div className="rounded-3xl border border-red-500/20 bg-red-500/10 px-6 py-12 text-center">
            <p className="text-lg font-medium text-white">
              Не удалось загрузить билеты
            </p>
            <p className="mt-2 text-sm text-red-200/80">
              {error instanceof Error
                ? error.message
                : "Ошибка запроса к серверу"}
            </p>
          </div>
        ) : !tickets || tickets.length === 0 ? (
          <div className="rounded-3xl border border-white/8 bg-white/[0.03] px-6 py-16 text-center">
            <p className="text-lg font-medium text-white">Билетов пока нет</p>
            <p className="mt-2 text-sm text-gray-400">
              Когда вы оплатите бронирование, билеты появятся здесь.
            </p>
            <button
              onClick={() => navigate("/")}
              className="mt-6 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-blue-500"
            >
              Выбрать фильм
            </button>
          </div>
        ) : (
          <div className="space-y-10">
            <section className="grid gap-4 sm:grid-cols-3">
              <div className="rounded-3xl border border-white/8 bg-white/[0.03] p-5">
                <p className="text-sm text-gray-400">Активные</p>
                <p className="mt-2 text-3xl font-semibold text-emerald-300">
                  {activeTickets.length}
                </p>
              </div>
              <div className="rounded-3xl border border-white/8 bg-white/[0.03] p-5">
                <p className="text-sm text-gray-400">Отменённые</p>
                <p className="mt-2 text-3xl font-semibold text-amber-300">
                  {deactivatedTickets.length}
                </p>
              </div>
              <div className="rounded-3xl border border-white/8 bg-white/[0.03] p-5">
                <p className="text-sm text-gray-400">Использованные</p>
                <p className="mt-2 text-3xl font-semibold text-slate-300">
                  {usedTickets.length}
                </p>
              </div>
            </section>

            <div className="space-y-5">
              <Section
                title="Активные билеты"
                description="Билеты, которые уже оплачены и готовы к посещению сеанса."
                tickets={activeTickets}
                accent="bg-emerald-400"
                defaultOpen
              />
              <Section
                title="Отменённые билеты"
                description="Билеты по отменённым бронированиям."
                tickets={deactivatedTickets}
                accent="bg-amber-400"
              />
              <Section
                title="Использованные билеты"
                description="Билеты, срок действия которых уже истёк."
                tickets={usedTickets}
                accent="bg-slate-300"
              />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
