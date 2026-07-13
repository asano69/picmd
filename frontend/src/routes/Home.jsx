import NavBar from "../components/NavBar";
import pb from "../lib/pb";

export default function Home() {
  return (
    <div class="mx-auto flex min-h-screen w-full max-w-xl flex-col items-center bg-[var(--color-bg)] px-6 py-12 text-[var(--color-text)]">
      <NavBar />
    </div>
  );
}
