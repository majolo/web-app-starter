// Note: this file is not generated
import { Api } from "@/gen/apis";
import { createClientComponentClient } from "@supabase/auth-helpers-nextjs";

export function GetAPI() {
  const supabase = createClientComponentClient();
  const handleSignUp = async () => {
    const { session } = await supabase.auth.getSession();
    return session;
  };
  handleSignUp().then((r) => console.log(r));

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
