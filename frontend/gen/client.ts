// Note: this file is not generated
import { Api } from "@/gen/apis";

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
