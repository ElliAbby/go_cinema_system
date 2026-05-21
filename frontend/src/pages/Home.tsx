import React from "react";
import { Link } from "react-router-dom";
import { useMovies } from "../hooks/useMovies";
import { MovieCard } from "../components/MovieCard";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const Home: React.FC = () => {
  const { data: movies, isLoading, error } = useMovies();

  return (
    <div className="min-h-screen bg-dark-950">
      {/* Hero */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-700 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-2xl">
            <h1 className="text-5xl font-bold text-white mb-4">
              Бронируйте билеты онлайн
            </h1>
            <p className="text-xl text-blue-100 mb-8">
              Выбирайте любимые фильмы и места в кинотеатре прямо из дома
            </p>
            <div className="flex gap-4">
              <Link
                to="/cinemas"
                className="px-6 py-3 bg-white text-blue-600 rounded-lg font-semibold hover:bg-blue-50 transition"
              >
                Выбрать кинотеатр
              </Link>
              <Link
                to="/#movies"
                className="px-6 py-3 bg-blue-700 text-white rounded-lg font-semibold hover:bg-blue-800 transition"
              >
                Посмотреть фильмы
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Movies Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16" id="movies">
        <h2 className="text-3xl font-bold text-white mb-8">Сейчас в кино</h2>

        {isLoading ? (
          <LoadingSpinner />
        ) : error ? (
          <div className="text-center py-12">
            <p className="text-red-400">Ошибка при загрузке фильмов</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {movies?.map((movie) => (
              <MovieCard key={movie.id} movie={movie} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
