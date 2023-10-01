import Link from 'next/link'

export const dynamic = 'force-dynamic'

export default async function Index() {
  return (
    <div className="w-full flex flex-col items-center">
      <div className="py-10">
        Hello there friend 👋
        <Link
          href="/diary"
          className="py-2 px-3 flex rounded-md no-underline bg-btn-background hover:bg-btn-background-hover"
        >
          Diary
        </Link>
      </div>
    </div>
  )
}
