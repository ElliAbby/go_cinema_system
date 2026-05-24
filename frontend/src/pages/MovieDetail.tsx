import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useMovie, useMovieSessions } from "../hooks/useMovies";
import { useCinema, useHall } from "../hooks/useCinemas";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { getMoviePoster } from "../utils/moviePoster";

export const MovieDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const movieId = id ? parseInt(id) : null;

  const { data: movie, isLoading: movieLoading } = useMovie(movieId);
  const { data: sessions, isLoading: sessionsLoading } =
    useMovieSessions(movieId);

  if (movieLoading) return <FullPageLoader />;

  if (!movie) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Фильм не найден</p>
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

  const poster = getMoviePoster(movie.title);

  return (
    <div className="min-h-screen bg-dark-950">
      {/* Header */}
      <div className="bg-gradient-to-br from-blue-600 via-indigo-700 to-purple-800 py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col gap-8 md:flex-row md:items-end">
          <div className="flex-shrink-0">
            <div className="relative w-48 overflow-hidden rounded-2xl border border-white/10 shadow-2xl shadow-black/30">
              <img
                src={poster}
                alt={movie.title}
                className="h-full w-full object-cover"
              />
            </div>
          </div>
          <div>
            <h1 className="text-4xl font-bold text-white mb-4">
              {movie.title}
            </h1>
            <div className="space-y-2 text-blue-100">
              <p>Продолжительность: {movie.duration} минут</p>
              <p>Рейтинг: {movie.rating}</p>
              <p>{movie.description}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Sessions */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h2 className="text-3xl font-bold text-white mb-8">Доступные сеансы</h2>

        {sessionsLoading ? (
          <LoadingSpinner />
        ) : !sessions || sessions.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-400">
              Нет доступных сеансов для этого фильма
            </p>
          </div>
        ) : (
          <div className="grid gap-4">
            {sessions.map((session) => (
              <SessionCard
                key={session.id}
                session={session}
                onBook={() => navigate(`/booking/${session.id}`)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

const SessionCard: React.FC<{
  session: {
    id: number;
    hall_id: number;
    start_time: string;
    price_base: number;
    available_seats?: number;
  };
  onBook: () => void;
}> = ({ session, onBook }) => {
  const { data: hall } = useHall(session.hall_id);
  const { data: cinema } = useCinema(hall?.cinema_id ?? null);

  return (
    <div className="bg-dark-800 border border-dark-700 rounded-lg p-6 hover:border-dark-600 transition">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div className="space-y-2">
          <p className="text-lg font-semibold text-white">
            {new Date(session.start_time).toLocaleTimeString("ru-RU", {
              hour: "2-digit",
              minute: "2-digit",
            })}
          </p>
          <p className="text-sm text-gray-400">
            {new Date(session.start_time).toLocaleDateString("ru-RU")}
          </p>
          <p className="text-sm text-gray-500">
            {cinema?.name ?? "Кинотеатр"}
            {hall?.name ? ` · ${hall.name}` : ""}
          </p>
          <p className="text-xs text-gray-500">{cinema?.address ?? ""}</p>
        </div>
        <div className="space-y-2 text-right">
          <p className="text-2xl font-bold text-blue-400">
            {Math.round(session.price_base)} ₽
          </p>
          <p className="text-sm text-gray-400">Цена билета</p>
        </div>
        <button
          onClick={onBook}
          disabled={session.available_seats === 0}
          className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white font-medium rounded-lg transition"
        >
          Выбрать места
        </button>
      </div>
    </div>
  );
};

const FullPageLoader = () => (
  <div className="min-h-screen bg-dark-950 flex items-center justify-center">
    <LoadingSpinner />
  </div>
);
