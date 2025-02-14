import axios from "axios";
axios.defaults.withCredentials = true;
const api = axios.create({
  baseURL: "http://localhost:8080",
  withCredentials: true,
});

export const getAllTierlists = () => {
  return api.get("/tierlist/");
};

export const whoAmI = () => {
  return api.get("/auth/me");
};
