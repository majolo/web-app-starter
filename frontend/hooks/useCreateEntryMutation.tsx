import { useMutation, useQueryClient } from "@tanstack/react-query";
import { GetAPI } from "@/gen/client";
import { DiaryV1CreateEntryRequest } from "@/gen/apis";

const useCreateEntryMutation = () => {
  const queryClient = useQueryClient();
  return useMutation(
    ["CreateEntryMutation"],
    (entry: DiaryV1CreateEntryRequest) =>
      GetAPI()
        .diaryServiceCreateEntry(entry)
        .then((res) => res.data)
        .catch((err) => err),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(["ListEntriesQuery"]);
      },
    },
  );
};

export default useCreateEntryMutation;
