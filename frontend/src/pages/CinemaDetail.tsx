import React from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useCinema, useHalls } from "../hooks/useCinemas";
import { LoadingSpinner } from "../components/LoadingSpinner";
import { getCinemaPoster } from "../utils/cinemaPoster";

export const CinemaDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const cinemaId = id ? parseInt(id) : null;

  const { data: cinema, isLoading: cinemaLoading } = useCinema(cinemaId);
  const { data: halls, isLoading: hallsLoading } = useHalls(cinemaId);

  if (cinemaLoading) return <FullPageLoader />;

  if (!cinema) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-400 mb-4">Кинотеатр не найден</p>
          <button
            onClick={() => navigate("/cinemas")}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
          >
            Вернуться к кинотеатрам
          </button>
        </div>
      </div>
    );
  }

  const poster = getCinemaPoster();

  return (
    <div className="min-h-screen bg-dark-950">
      {/* Header */}
      <div className="bg-gradient-to-br from-purple-600 via-fuchsia-700 to-pink-700 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col gap-8 md:flex-row md:items-end">
          <div className="flex-shrink-0">
            <div className="relative w-56 overflow-hidden rounded-2xl border border-white/10 shadow-2xl shadow-black/30">
              <img
                src={poster}
                alt={cinema.name}
                className="h-full w-full object-cover"
              />
            </div>
          </div>
          <div>
            <h1 className="text-4xl font-bold text-white mb-2">
              {cinema.name}
            </h1>
            <p className="text-purple-100">📍 {cinema.address}</p>
          </div>
        </div>
      </div>

      {/* Halls */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h2 className="text-3xl font-bold text-white mb-8">Залы</h2>

        {hallsLoading ? (
          <LoadingSpinner />
        ) : !halls || halls.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-400">Нет доступных залов</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {halls.map((hall) => (
              <div
                key={hall.id}
                className="bg-dark-800 border border-dark-700 rounded-lg p-6 hover:border-dark-600 transition cursor-pointer"
                onClick={() => navigate(`/cinema/${cinemaId}/hall/${hall.id}`)}
              >
                <h3 className="text-xl font-semibold text-white mb-4">
                  {hall.name}
                </h3>
                <button className="w-full mt-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition">
                  Посмотреть сеансы
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

const FullPageLoader = () => (
  <div className="min-h-screen bg-dark-950 flex items-center justify-center">
    <LoadingSpinner />
  </div>
);
