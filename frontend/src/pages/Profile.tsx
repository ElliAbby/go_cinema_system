import React from "react";
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
  const displayName = profile.email?.split("@")[0] || "Пользователь";

  return (
    <div className="min-h-screen bg-dark-950">
      <div className="mx-auto max-w-4xl px-4 py-10 sm:px-6 lg:px-8 lg:py-14">
        <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.24em] text-gray-500">
              Account
            </p>
            <h1 className="mt-2 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
              Мой профиль
            </h1>
            <p className="mt-2 max-w-2xl text-sm text-gray-400 sm:text-base">
              Аккаунт, контактные данные и статус доступа.
            </p>
          </div>

          <div
            className={`inline-flex w-fit items-center gap-2 rounded-full border px-3.5 py-2 text-sm font-medium ${accountStatus.className}`}
          >
            <span
              className={`h-2 w-2 rounded-full ${profile.is_active === false ? "bg-red-400" : "bg-emerald-400"}`}
            />
            {accountStatus.label}
          </div>
        </div>

        <section className="overflow-hidden rounded-3xl border border-white/8 bg-white/[0.03] shadow-[0_24px_60px_rgba(0,0,0,0.35)] backdrop-blur">
          <div className="border-b border-white/8 px-6 py-6 sm:px-8 sm:py-8">
            <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-2xl border border-white/10 bg-gradient-to-br from-white/10 to-white/[0.03] text-2xl font-semibold text-white">
                {avatarLetter}
              </div>

              <div className="min-w-0 flex-1">
                <h2 className="truncate text-2xl font-semibold text-white sm:text-[28px]">
                  {displayName}
                </h2>
                <p className="mt-1 break-all text-sm text-gray-400 sm:text-base">
                  {profile.email}
                </p>
              </div>
            </div>
          </div>

          <div className="grid gap-4 p-6 sm:p-8 md:grid-cols-3">
            <div className="rounded-2xl border border-white/8 bg-black/20 p-5">
              <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
                Телефон
              </p>
              <p className="mt-3 text-base font-medium text-white">
                {profile.phone || "Не указан"}
              </p>
            </div>

            <div className="rounded-2xl border border-white/8 bg-black/20 p-5">
              <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
                Регистрация
              </p>
              <p className="mt-3 text-base font-medium text-white">
                {createdAt}
              </p>
            </div>

            <div className="rounded-2xl border border-white/8 bg-black/20 p-5">
              <p className="text-xs uppercase tracking-[0.2em] text-gray-500">
                Обновлен
              </p>
              <p className="mt-3 text-base font-medium text-white">
                {updatedAt}
              </p>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
};
