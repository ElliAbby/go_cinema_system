import React from "react";
import { useAuth } from "../context/AuthContext";
import { useProfile } from "../hooks/useProfile";
import { LoadingSpinner } from "../components/LoadingSpinner";

const getAccountStatus = (isActive?: boolean) => {
  if (isActive === false) {
    return {
      label: "Неактивен",
      className: "bg-red-900/30 text-red-300 border-red-800/70",
    };
  }

  return {
    label: "Активен",
    className: "bg-emerald-900/30 text-emerald-300 border-emerald-800/70",
  };
};

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
  const accountStatus = getAccountStatus(profile.is_active);
  const avatarLetter = profile.email?.trim().charAt(0).toUpperCase() || "U";
  const displayName =
    profile.email?.split("@")[0] || `Пользователь #${profile.id}`;

  return (
    <div className="min-h-screen bg-dark-950 relative overflow-hidden">
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute -top-32 right-[-6rem] h-72 w-72 rounded-full bg-cyan-500/10 blur-3xl" />
        <div className="absolute top-40 left-[-8rem] h-80 w-80 rounded-full bg-amber-500/10 blur-3xl" />
      </div>

      <div className="relative max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-10 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-4xl font-bold text-white tracking-tight">
              Мой профиль
            </h1>
            <p className="mt-2 text-gray-400 max-w-2xl">
              Личные данные аккаунта и текущий статус профиля.
            </p>
          </div>
          <div className="inline-flex w-fit items-center gap-2 rounded-full border border-dark-700 bg-dark-800/80 px-4 py-2 text-sm text-gray-300 shadow-lg shadow-black/20 backdrop-blur">
            <span
              className={`h-2.5 w-2.5 rounded-full ${profile.is_active === false ? "bg-red-400" : "bg-emerald-400"}`}
            />
            {accountStatus.label}
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[1.15fr_0.85fr]">
          <section className="relative overflow-hidden rounded-3xl border border-dark-700 bg-gradient-to-br from-dark-800 via-dark-800 to-dark-900 shadow-2xl shadow-black/30">
            <div className="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-cyan-400 via-amber-400 to-emerald-400" />
            <div className="p-6 sm:p-8">
              <div className="flex flex-col gap-6 sm:flex-row sm:items-center">
                <div className="flex h-20 w-20 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-500 to-amber-500 text-3xl font-bold text-white shadow-lg shadow-cyan-500/20">
                  {avatarLetter}
                </div>

                <div className="min-w-0 flex-1">
                  <h2 className="text-2xl font-semibold text-white break-words">
                    {displayName}
                  </h2>
                  <p className="mt-1 text-gray-400 break-all">
                    {profile.email}
                  </p>
                </div>
              </div>

              <div className="mt-8 grid gap-4 sm:grid-cols-3">
                <div className="rounded-2xl border border-dark-700 bg-black/20 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-gray-500">
                    Статус
                  </p>
                  <div
                    className={`mt-3 inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium ${accountStatus.className}`}
                  >
                    {accountStatus.label}
                  </div>
                </div>

                <div className="rounded-2xl border border-dark-700 bg-black/20 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-gray-500">
                    Телефон
                  </p>
                  <p className="mt-3 text-white font-medium">
                    {profile.phone || "Не указан"}
                  </p>
                </div>

                <div className="rounded-2xl border border-dark-700 bg-black/20 p-4">
                  <p className="text-xs uppercase tracking-[0.18em] text-gray-500">
                    Дата регистрации
                  </p>
                  <p className="mt-3 text-white font-medium">{createdAt}</p>
                </div>
              </div>
            </div>
          </section>
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
