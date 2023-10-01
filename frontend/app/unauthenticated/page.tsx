import {createServerComponentClient} from "@supabase/auth-helpers-nextjs";
import {cookies} from "next/headers";
import {redirect} from "next/navigation";

export default async function Unauthenticated() {
  const supabase = createServerComponentClient({cookies})
  const {data:{session}} = await supabase.auth.getSession()

  if (session) {
    redirect('/')
  }

  return <>You are unauthenticated, please login.</>
}
