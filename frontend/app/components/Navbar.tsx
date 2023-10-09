import { createServerComponentClient } from "@supabase/auth-helpers-nextjs";
import { cookies } from "next/headers";
import Link from "next/link";
import LogoutButton from "./LogoutButton";

export const dynamic = "force-dynamic";

export default async function Navbar() {
  const supabase = createServerComponentClient({ cookies });
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <nav className="w-full flex justify-center h-16 bg-gray-400">
      <div />
      {user ? (
        <div className="flex items-center gap-4">
          Hey, {user.email}!
          <LogoutButton />
        </div>
      ) : (
        <Link
          href="/login"
          className="py-2 px-3 flex rounded-md no-underline bg-btn-background hover:bg-btn-background-hover"
        >
          Login
        </Link>
      )}
    </nav>
  );
}
