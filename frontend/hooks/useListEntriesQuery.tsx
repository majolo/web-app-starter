import { useQuery } from "@tanstack/react-query";
import { GetAPI } from "@/gen/client";
import { DiaryV1Entry } from "@/gen/apis";

const useListEntriesQuery = () =>
  useQuery({
    queryKey: ["ListEntriesQuery"],
    queryFn: (): Promise<DiaryV1Entry[]> =>
      GetAPI()
        .diaryServiceListEntries()
        .then((res) => res.data.entries as DiaryV1Entry[]),
  });

export default useListEntriesQuery;
