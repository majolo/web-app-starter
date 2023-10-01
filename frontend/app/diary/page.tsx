import Protected from "@/components/Protected";

export default async function Page() {
  return (
    <div className={"py-24"}>
      <Protected />
      Hello World.
    </div>
  );
}
