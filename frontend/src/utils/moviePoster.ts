const DEFAULT_POSTER = "/images/default.avif";

const POSTER_BY_TITLE: Record<string, string> = {
  интерстеллар: "/images/intestellar.png",
  "темный рыцарь": "/images/batman.png",
  "тёмный рыцарь": "/images/batman.png",
  начало: "/images/start.png",
  матрица: "/images/matrix.png",
  аватар: "/images/avatar.png",
  оппенгеймер: "/images/opengamer.jpg",
  опенгеймер: "/images/opengamer.jpg",
  "дюна: часть вторая": "/images/dune.jpg",
  "дюна 2": "/images/dune.jpg",
  дюна: "/images/dune.jpg",
  "человек-паук: нет пути домой": "/images/spider_man.webp",
  "человек паук: нет пути домой": "/images/spider_man.webp",
  "человек-паук": "/images/spider_man.webp",
  "джон уик 4": "/images/jhon_week.png",
  "джон уик": "/images/jhon_week.png",
  гладиатор: "/images/gladiator.jpg",
  "гладиатор 2": "/images/gladiator.jpg",
  фуриоса: "/images/furiosa.jpg",
  "дэдпул и росомаха": "/images/deadpool.webp",
};

export const getMoviePoster = (title: string) => {
  const key = title.trim().toLowerCase();

  return POSTER_BY_TITLE[key] || DEFAULT_POSTER;
};
