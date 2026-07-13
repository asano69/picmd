import { createSignal } from "solid-js";
import { A } from "@solidjs/router";
import pb from "../lib/pb";

export default function NavBar(props) {
  const [refreshing, setRefreshing] = createSignal(false);

 
  const handleLogout = () => pb.authStore.clear();

  return (
    <div class="mb-10 flex w-full flex-wrap items-center justify-between gap-y-3">
      <A href="/" class="font-serif text-4xl flex items-center gap-2 transition-opacity hover:opacity-80">
        <img src="/favicon.svg" alt="" class="h-12 w-12" />
        <h1>myapp</h1>
      </A>
      <nav class="flex flex-wrap items-center gap-3">
      
        <A href="/stats" class="btn">Stats</A>
        <A href="/settings" class="btn">Settings</A>
        <button type="button" class="btn" onClick={handleLogout}>Log out</button>
      </nav>
    </div>
  );
}
