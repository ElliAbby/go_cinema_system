import React from "react";
import { Link } from "react-router-dom";
import { Cinema } from "../types";
import { getCinemaPoster } from "../utils/cinemaPoster";

interface CinemaCardProps {
  cinema: Cinema;
}

export const CinemaCard: React.FC<CinemaCardProps> = ({ cinema }) => {
  const poster = getCinemaPoster();

  return (
    <Link to={`/cinema/${cinema.id}`}>
      <div className="group cursor-pointer">
        <div className="relative overflow-hidden rounded-lg aspect-[4/3] mb-4 bg-dark-800 border border-dark-700 hover:shadow-lg hover:shadow-purple-500/20 transition-all duration-300">
          <img
            src={poster}
            alt={cinema.name}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/75 via-black/10 to-transparent" />
          <div className="absolute bottom-3 left-3 right-3">
            <span className="inline-flex rounded-full bg-black/60 px-2 py-1 text-[11px] font-medium text-white backdrop-blur-sm">
              Кинотеатр
            </span>
          </div>
        </div>
        <h3 className="text-lg font-semibold text-white group-hover:text-purple-400 transition mb-2">
          {cinema.name}
        </h3>
        <p className="text-sm text-gray-400">{cinema.address}</p>
      </div>
    </Link>
  );
};
