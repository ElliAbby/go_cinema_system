import React from "react";
import { useAuth } from "../context/AuthContext";
import { useProfile } from "../hooks/useProfile";
import { LoadingSpinner } from "../components/LoadingSpinner";

export const Profile: React.FC = () => {
  const { user } = useAuth();
  const { data: profile, isLoading, isError } = useProfile();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-dark-950 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (isError || !profile) {
    return (
      <div className="min-h-screen bg-dark-950">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <h1 className="text-4xl font-bold text-white mb-2">Мой профиль</h1>
          <p className="text-gray-400 mb-8">Не удалось загрузить профиль</p>
          <div className="bg-red-900/20 border border-red-700 rounded-lg p-4 text-red-300 text-sm">
            Проверьте авторизацию и попробуйте обновить страницу.
          </div>
        </div>
      </div>
    );
  }

  const createdAt = profile.created_at
    ? new Date(profile.created_at).toLocaleString("ru-RU")
    : "Не указано";
  const updatedAt = profile.updated_at
    ? new Date(profile.updated_at).toLocaleString("ru-RU")
    : "Не указано";

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-4xl font-bold text-white mb-2">Мой профиль</h1>
        <p className="text-gray-400 mb-10">Ваши неприватные данные аккаунта</p>

        <div className="bg-dark-800 border border-dark-700 rounded-xl overflow-hidden">
          <div className="px-6 py-4 border-b border-dark-700">
            <p className="text-lg text-white font-semibold">
              Пользователь #{profile.id}
            </p>
          </div>

          <div className="divide-y divide-dark-700">
            <div className="px-6 py-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <span className="text-gray-400">Email</span>
              <span className="sm:col-span-2 text-white break-all">
                {profile.email}
              </span>
            </div>

            <div className="px-6 py-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <span className="text-gray-400">Телефон</span>
              <span className="sm:col-span-2 text-white">
                {profile.phone || "Не указан"}
              </span>
            </div>

            <div className="px-6 py-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <span className="text-gray-400">Статус аккаунта</span>
              <span className="sm:col-span-2 text-white">
                {profile.is_active === false ? "Неактивен" : "Активен"}
              </span>
            </div>

            <div className="px-6 py-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <span className="text-gray-400">Дата регистрации</span>
              <span className="sm:col-span-2 text-white">{createdAt}</span>
            </div>

            <div className="px-6 py-4 grid grid-cols-1 sm:grid-cols-3 gap-3">
              <span className="text-gray-400">Последнее обновление</span>
              <span className="sm:col-span-2 text-white">{updatedAt}</span>
            </div>
          </div>
        </div>

        {user && user.id !== profile.id && (
          <p className="mt-6 text-sm text-yellow-400">
            Профиль загружен по токену, локальный user id отличается.
          </p>
        )}
      </div>
    </div>
  );
};
