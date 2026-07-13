// frontend/src/main.jsx
import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { createSignal, onMount, onCleanup, Show } from "solid-js";

// Order matters: tokens.css defines the CSS custom properties every other
// stylesheet consumes via var().
import "./style.css";
import Settings from "./routes/Settings";
import Login from "./routes/Login";
import Upload from "./routes/Upload";
import pb from "./lib/pb";

// AuthGate blocks the whole app behind Login until a valid superuser
// session exists, tracking pb.authStore so it reacts immediately to
// both login and logout.
//
// pb.authStore.isValid only checks the token locally (decoded JWT expiry),
// so it can still report "valid" right after the tab wakes up from being
// suspended, or if the session was revoked server-side in the meantime.
// To catch that, the session is re-verified against the server whenever
// the tab becomes visible again; a rejected refresh clears the store,
// which flips the app back to Login via the onChange subscription above.
function AuthGate(props) {
  const [authed, setAuthed] = createSignal(pb.authStore.isValid);
  const unsubscribe = pb.authStore.onChange(() =>
    setAuthed(pb.authStore.isValid),
  );
  onCleanup(unsubscribe);

  const verifySession = async () => {
    if (!pb.authStore.isValid) return;
    try {
      await pb.collection("_superusers").authRefresh();
    } catch {
      pb.authStore.clear();
    }
  };

  const handleVisibility = () => {
    if (document.visibilityState === "visible") verifySession();
  };

  onMount(() => {
    verifySession();
    document.addEventListener("visibilitychange", handleVisibility);
  });
  onCleanup(() =>
    document.removeEventListener("visibilitychange", handleVisibility),
  );

  return (
    <Show when={authed()} fallback={<Login />}>
      {props.children}
    </Show>
  );
}

render(
  () => (
    <AuthGate>
      <Router>
        <Route path="/" component={Upload} />
        <Route path="/settings" component={Settings} />
      </Router>
    </AuthGate>
  ),
  document.getElementById("app"),
);
