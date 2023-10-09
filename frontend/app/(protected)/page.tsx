"use client";
import useListEntriesQuery from "@/hooks/useListEntriesQuery";
import { useState } from "react";
import { useCreateEntryMutation } from "@/hooks";

export default function Page() {
  const { data: entries } = useListEntriesQuery();
  return (
    <div className={"py-24"}>
      <h1 className={"text-4xl font-bold text-center"}>Diary</h1>
      <TextBoxWithButton />
      {entries &&
        entries?.map((entry) => (
          <div key={entry.id}>
            <p className={"text-lg"}>{entry.text}</p>
            <p className={"text-lg"}>{entry.id}</p>
          </div>
        ))}
    </div>
  );
}

const TextBoxWithButton = () => {
  const [text, setText] = useState("");
  const { mutate } = useCreateEntryMutation();

  const handleSubmit = () => {
    mutate({ text });
  };

  return (
    <div className="flex flex-col items-center justify-center py-5">
      <div className="bg-gray-400 p-4 rounded-lg shadow-md ">
        <input
          type="text"
          value={text}
          onChange={(e) => setText(e.target.value)}
          className="border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
        <button
          onClick={handleSubmit}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline mt-2"
        >
          Submit
        </button>
      </div>
    </div>
  );
};
