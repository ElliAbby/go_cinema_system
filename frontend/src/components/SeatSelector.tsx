import React from "react";

interface SeatSelectorProps {
  rows: number;
  seatsPerRow: number;
  selectedSeats: Set<number>;
  onSeatClick: (seatNumber: number, seatId?: number) => void;
  seatIdMap?: Record<number, number>;
  reservedSeatIds?: Set<number>;
}

export const SeatSelector: React.FC<SeatSelectorProps> = ({
  rows,
  seatsPerRow,
  selectedSeats,
  onSeatClick,
  seatIdMap,
  reservedSeatIds,
}) => {
  const getSeatNumber = (row: number, seat: number) =>
    row * seatsPerRow + seat + 1;

  return (
    <div className="bg-dark-800 rounded-lg p-6">
      <div className="mb-8">
        <div className="h-1 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full mx-auto w-32 mb-4"></div>
        <p className="text-center text-sm text-gray-400">Экран</p>
      </div>

      <div className="flex justify-center">
        <div className="space-y-2">
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <div key={rowIndex} className="flex gap-2 justify-center">
              <span className="w-6 text-xs text-gray-500 flex items-center justify-center">
                {String.fromCharCode(65 + rowIndex)}
              </span>
              {Array.from({ length: seatsPerRow }).map((_, seatIndex) => {
                const seatNumber = getSeatNumber(rowIndex, seatIndex);
                const seatId = seatIdMap ? seatIdMap[seatNumber] : undefined;
                const isSelected = selectedSeats.has(seatId ?? seatNumber);
                const isReserved =
                  seatId !== undefined && reservedSeatIds?.has(seatId);

                return (
                  <button
                    key={seatIndex}
                    onClick={() => onSeatClick(seatNumber, seatId)}
                    disabled={isReserved}
                    className={`w-8 h-8 rounded text-xs font-medium transition-all duration-200 ${
                      isReserved
                        ? "bg-red-700 text-white cursor-not-allowed"
                        : isSelected
                          ? "bg-blue-500 text-white shadow-lg shadow-blue-500/50"
                          : "bg-dark-700 text-gray-400 hover:bg-dark-600"
                    }`}
                  >
                    {seatIndex + 1}
                  </button>
                );
              })}
            </div>
          ))}
        </div>
      </div>

      <div className="mt-8 flex justify-center gap-6">
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 bg-dark-700 rounded"></div>
          <span className="text-sm text-gray-400">Доступно</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 bg-blue-500 rounded"></div>
          <span className="text-sm text-gray-400">Выбрано</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 bg-red-700 rounded"></div>
          <span className="text-sm text-gray-400">Занято</span>
        </div>
      </div>
    </div>
  );
};
