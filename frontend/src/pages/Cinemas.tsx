import React from "react";
import { useCinemas } from "../hooks/useCinemas";
import { CinemaCard } from "../components/CinemaCard";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const Cinemas: React.FC = () => {
  const { data: cinemas, isLoading, error } = useCinemas();

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-4xl font-bold text-white mb-2">Кинотеатры</h1>
        <p className="text-gray-400 mb-12">Выберите ближайший кинотеатр</p>

        {isLoading ? (
          <LoadingSpinner />
        ) : error ? (
          <div className="text-center py-12">
            <p className="text-red-400">Ошибка при загрузке кинотеатров</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {cinemas?.map((cinema) => (
              <CinemaCard key={cinema.id} cinema={cinema} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
