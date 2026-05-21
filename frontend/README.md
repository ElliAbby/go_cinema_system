# Cinema App - Cinema Booking Frontend

Современный минималистичный frontend для системы бронирования билетов в кинотеатр с темной темой.
s

## 🎬 Особенности

- ✨ Современный и минималистичный дизайн
- 🌙 Полная темная тема
- ⚡ Быстрая загрузка (Vite)
- 📱 Адаптивный дизайн (мобильный, планшет, десктоп)
- 🔐 Аутентификация пользователей
- 🎫 Бронирование билетов и выбор мест
- 💳 Интеграция с платежной системой (Kafka)
- 📋 Управление бронированиями
- 🎨 Tailwind CSS для стилизации

## 🚀 Быстрый старт

### Требования

- Node.js 16+
- npm или yarn

### Установка зависимостей

```bash
cd frontend
npm install
```

### Разработка

```bash
npm run dev
```

Откройте [http://localhost:3000](http://localhost:3000) в браузере.

### Сборка

```bash
npm run build
```

Оптимизированная сборка будет создана в папке `dist`.

### Preview

```bash
npm run preview
```

## 📁 Структура проекта

```
frontend/
├── src/
│   ├── api/           # API клиенты
│   ├── components/    # Переиспользуемые компоненты
│   ├── context/       # React Context (Auth)
│   ├── hooks/         # Кастомные хуки для API
│   ├── pages/         # Страницы приложения
│   ├── types/         # TypeScript типы
│   ├── App.tsx        # Главный компонент
│   ├── main.tsx       # Точка входа
│   └── index.css      # Глобальные стили
├── index.html         # HTML шаблон
├── package.json       # Зависимости
├── tsconfig.json      # TypeScript конфиг
├── tailwind.config.js # Tailwind конфиг
└── vite.config.ts     # Vite конфиг
```

## 🔌 API Интеграция

Frontend подключается к API на `http://localhost:8080`:

- Автоматический прокси для `/api` на `localhost:8080`
- Токен JWT сохраняется в `localStorage`
- Автоматическое добавление токена в заголовок `Authorization`

### Переменные окружения

Создайте файл `.env` в папке `frontend`:

```env
VITE_API_BASE_URL=http://localhost:8080
```

## 🎨 Темная тема

Тема использует пользовательские цвета из Tailwind:

- `dark-50` до `dark-950` для оттенков серого
- Синий и фиолетовый градиенты для акцентов
- Гладкие переходы и hover эффекты

## 📦 Основные технологии

- **React 18** - UI библиотека
- **TypeScript** - Типизация
- **Vite** - Build tool
- **Tailwind CSS** - Стилизация
- **React Query** - Data fetching
- **React Router** - Маршрутизация
- **Axios** - HTTP клиент
- **date-fns** - Работа с датами

## 🔐 Аутентификация

- Регистрация новых пользователей
- Вход по email/пароль
- JWT токены с временем жизни 1 час
- Автоматическое перенаправление на login при истечении сессии

## 📝 Доступные страницы

- `/` - Главная (список фильмов)
- `/cinemas` - Список кинотеатров
- `/cinema/:id` - Детали кинотеатра и залы
- `/movie/:id` - Детали фильма и сеансы
- `/booking/:sessionId` - Выбор мест
- `/booking-confirm/:bookingId` - Подтверждение и оплата
- `/bookings` - Мои бронирования
- `/login` - Вход
- `/register` - Регистрация

## 🛠️ Лоучинг и развитие

Для локальной разработки:

```bash
# Запустить backend на 8080
cd ..
go run cmd/cinema-service/main.go

# Запустить frontend на 3000 (в новом терминале)
cd frontend
npm run dev
```

## 📄 Лицензия

MIT
