"use client";
import useListEntriesQuery from "@/hooks/useListEntriesQuery";

export default function Page() {
  const { data: entries } = useListEntriesQuery();
  return (
    <div className={"py-24"}>
      <h1 className={"text-4xl font-bold text-center"}>Diary</h1>
      {entries &&
        entries?.map((entry) => (
          <div key={entry.id}>
            <p className={"text-lg"}>{entry.text}</p>
          </div>
        ))}
    </div>
  );
}
