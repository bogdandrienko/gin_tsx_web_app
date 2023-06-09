// TODO export /////////////////////////////////////////////////////////////////////////////////////////////////////////

// export const DEBUG_CONSTANT = process.env.NODE_ENV === "production";
export const DEBUG_CONSTANT = true;

// export const SERVER_HOST_AND_PORT_CONSTANT = "http://127.0.0.1:8003"; // todo DEV
export const SERVER_HOST_AND_PORT_CONSTANT = "http://172.30.23.16:8002"; // todo PROD
// export const SERVER_HOST_AND_PORT_CONSTANT = "http://172.30.23.16:8003"; // todo TEST

export class HttpMethods {
  static GET() {
    return "GET";
  }
  static READ() {
    return "GET";
  }
  static POST() {
    return "POST";
  }
  static CREATE() {
    return "POST";
  }
  static PUT() {
    return "PUT";
  }
  static UPDATE() {
    return "PUT";
  }
  static DELETE() {
    return "DELETE";
  }
}

export const news1 = [
  {
    Title: "Личный профиль:",
    Status: "disable",
    Link: "#",
    Description:
      "расширение профиля: дополнительная информация: образование, хобби, интересы, изображение, статистика, " +
      "достижения и участие",
    Helps: "",
    Danger: "",
  },
  {
    Title: "Рационализаторские предложения:",
    Status: "disable",
    Link: "#",
    Description:
      "шаблон и подача рац. предложений, модерация и общий список с рейтингами и комментариями",
    Helps: "",
    Danger: "",
  },
  {
    Title: "Банк идей:",
    Status: "active",
    Link: "/idea_create",
    Description:
      "подача и редактирование, шаблон, модерация, комментирование, рейтинги, лучшие идеи и списки лидеров...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "Инструкции: видео и текстовые",
    Status: "active",
    Link: "/video_study",
    Description: "лента с информацией по веб-платформе",
    Helps: "материал будет своевременно обновляться",
    Danger: "",
  },
  {
    Title: "Отпуска:",
    Status: "active",
    Link: "/vacation",
    Description: "выгрузка Ваших данных по отпуску за выбранный период",
    Helps: "",
    Danger: "",
  },
  {
    Title: "Расчётный лист:",
    Status: "active",
    Link: "/salary",
    Description: "выгрузка Вашего расчётного листа за выбранный период",
    Helps: "",
    Danger: "'контрактникам' выгрузка недоступна!",
  },
];

export const news = [
  {
    Title: "1:",
    Status: "active",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "2:",
    Status: "disable",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "3:",
    Status: "disable",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "4:",
    Status: "disable",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "5:",
    Status: "disable",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
  {
    Title: "6:",
    Status: "disable",
    Link: "#",
    Description: "...",
    Helps: "",
    Danger: "",
  },
];
