import React from "react";
import { Link } from "react-router-dom";
import { Movie } from "../types";
import { getMoviePoster } from "../utils/moviePoster";

interface MovieCardProps {
  movie: Movie;
}

export const MovieCard: React.FC<MovieCardProps> = ({ movie }) => {
  const poster = getMoviePoster(movie.title);

  return (
    <Link to={`/movie/${movie.id}`}>
      <div className="group cursor-pointer h-full">
        <div className="relative rounded-lg aspect-[2/3] mb-4 overflow-hidden bg-dark-800 border border-dark-700 hover:shadow-lg hover:shadow-blue-500/20 transition-all duration-300">
          <img
            src={poster}
            alt={movie.title}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent" />
          <div className="absolute bottom-3 left-3 right-3">
            <span className="inline-flex rounded-full bg-black/60 px-2 py-1 text-[11px] font-medium text-white backdrop-blur-sm">
              Превью фильма
            </span>
          </div>
        </div>
        <h3 className="text-lg font-semibold text-white group-hover:text-blue-400 transition mb-1 line-clamp-2">
          {movie.title}
        </h3>
        <div className="flex items-center gap-2 text-sm text-gray-400 mb-2">
          <span>{movie.duration} мин</span>
          <span>•</span>
          <span className="px-2 py-0.5 bg-dark-700 rounded text-xs">
            {movie.rating}
          </span>
        </div>
        <p className="text-sm text-gray-400 line-clamp-2">
          {movie.description}
        </p>
      </div>
    </Link>
  );
};
