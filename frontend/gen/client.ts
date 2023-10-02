// Note: this file is not generated
import { Api } from "@/gen/apis";
import * as process from "process";

// TODO: do the auth here too
export function GetAPI() {
  const customFetch: typeof fetch = async (input, init) => {
    return fetch(input, {
      ...init,
      credentials: "include", // Automatically include credentials
    });
  };

  return new Api({
    baseUrl: "http://localhost:8080",
    customFetch: customFetch,
  }).api;
}
