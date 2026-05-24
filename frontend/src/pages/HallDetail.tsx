import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { useCinema, useHall, useHallSessions } from "../hooks/useCinemas";
import { useMovie } from "../hooks/useMovies";

const formatDateTime = (value?: string) => {
  if (!value) return "—";

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "long",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
};

const SessionCard: React.FC<{
  session: {
    id: number;
    movie_id: number;
    start_time: string;
    price_base: number;
  };
}> = ({ session }) => {
  const navigate = useNavigate();
  const { data: movie } = useMovie(session.movie_id);

  return (
    <div className="rounded-2xl border border-dark-700 bg-dark-800 p-5 transition hover:border-dark-600">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-lg font-semibold text-white">
            {movie?.title ?? "Фильм"}
          </p>
          <p className="mt-1 text-sm text-gray-400">
            {formatDateTime(session.start_time)}
          </p>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-blue-400">
            {Math.round(session.price_base)} ₽
          </p>
          <p className="text-sm text-gray-500">за билет</p>
        </div>
        <button
          onClick={() => navigate(`/booking/${session.id}`)}
          className="rounded-lg bg-blue-600 px-5 py-3 font-medium text-white transition hover:bg-blue-700"
        >
          Выбрать места
        </button>
      </div>
    </div>
  );
};

export const HallDetail: React.FC = () => {
  const { cinemaId, hallId } = useParams<{
    cinemaId: string;
    hallId: string;
  }>();
  const navigate = useNavigate();

  const cinemaIdValue = cinemaId ? parseInt(cinemaId, 10) : null;
  const hallIdValue = hallId ? parseInt(hallId, 10) : null;

  const { data: cinema, isLoading: cinemaLoading } = useCinema(cinemaIdValue);
  const { data: hall, isLoading: hallLoading } = useHall(hallIdValue);
  const { data: sessions, isLoading: sessionsLoading } =
    useHallSessions(hallIdValue);

  if (cinemaLoading || hallLoading) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (!cinema || !hall) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Кинотеатр или зал не найден</p>
          <button
            onClick={() => navigate("/cinemas")}
            className="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
          >
            Вернуться к кинотеатрам
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <button
          onClick={() => navigate(`/cinema/${cinema.id}`)}
          className="mb-8 text-blue-400 transition hover:text-blue-300"
        >
          ← Назад к кинотеатру
        </button>

        <div className="rounded-3xl border border-dark-700 bg-dark-900 p-8">
          <p className="text-sm uppercase tracking-[0.2em] text-gray-500">
            Зал
          </p>
          <h1 className="mt-2 text-4xl font-bold text-white">{hall.name}</h1>
          <p className="mt-2 text-gray-400">
            {cinema.name} · {cinema.address}
          </p>
        </div>

        <div className="mt-10">
          <h2 className="mb-6 text-3xl font-bold text-white">
            Сеансы в этом зале
          </h2>

          {sessionsLoading ? (
            <div className="flex min-h-[30vh] items-center justify-center">
              <LoadingSpinner />
            </div>
          ) : !sessions || sessions.length === 0 ? (
            <div className="rounded-2xl border border-dark-700 bg-dark-800 p-8 text-center text-gray-400">
              В этом зале пока нет сеансов
            </div>
          ) : (
            <div className="space-y-4">
              {sessions.map((session) => (
                <SessionCard key={session.id} session={session} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
