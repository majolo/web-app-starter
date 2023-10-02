"use client";
import useListEntriesQuery from "@/hooks/useListEntriesQuery";
import { DiaryV1Entry } from "@/gen/apis";

export default function Page() {
  const { data: entries, isLoading: isLoadingEntries } = useListEntriesQuery();
  console.log(entries);
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
