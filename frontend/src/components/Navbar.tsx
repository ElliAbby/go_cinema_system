import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export const Navbar: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <nav className="sticky top-0 z-50 bg-dark-900 border-b border-dark-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-lg">🎬</span>
            </div>
            <span className="text-xl font-bold text-white">Cinema App</span>
          </Link>

          {/* Menu */}
          <div className="flex items-center gap-6">
            <Link to="/" className="text-gray-300 hover:text-white transition">
              Главная
            </Link>
            <Link
              to="/cinemas"
              className="text-gray-300 hover:text-white transition"
            >
              Кинотеатры
            </Link>
            {isAuthenticated && (
              <Link
                to="/bookings"
                className="text-gray-300 hover:text-white transition"
              >
                Мои заказы
              </Link>
            )}
            {isAuthenticated && (
              <Link
                to="/profile"
                className="text-gray-300 hover:text-white transition"
              >
                Профиль
              </Link>
            )}
          </div>

          {/* Auth */}
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <>
                <div className="text-sm text-gray-400">{user?.email}</div>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 rounded-lg bg-dark-700 hover:bg-dark-600 text-gray-200 transition text-sm font-medium"
                >
                  Выход
                </button>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  className="px-4 py-2 rounded-lg bg-dark-700 hover:bg-dark-600 text-gray-200 transition text-sm font-medium"
                >
                  Вход
                </Link>
                <Link
                  to="/register"
                  className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white transition text-sm font-medium"
                >
                  Регистрация
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};
