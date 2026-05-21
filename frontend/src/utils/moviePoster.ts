const DEFAULT_POSTER = "/images/avatar.png";

const POSTER_BY_TITLE: Record<string, string> = {
  интерстеллар: "/images/intestellar.png",
  "темный рыцарь": "/images/batman.png",
  "тёмный рыцарь": "/images/batman.png",
  начало: "/images/start.png",
  матрица: "/images/matrix.png",
  аватар: "/images/avatar.png",
};

export const getMoviePoster = (title: string) => {
  const key = title.trim().toLowerCase();

  return POSTER_BY_TITLE[key] || DEFAULT_POSTER;
};
