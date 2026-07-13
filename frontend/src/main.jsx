import { render } from "solid-js/web";
import { Router, Route } from "@solidjs/router";
import { createSignal, onCleanup, Show } from "solid-js";

// Order matters: tokens.css defines the CSS custom properties every other
// stylesheet consumes via var().
import "./style.css";
import Home from "./routes/Home";
import Settings from "./routes/Settings";
import Login from "./routes/Login";

import pb from "./lib/pb";

// AuthGate blocks the whole app behind Login until a valid superuser
// session exists, tracking pb.authStore so it reacts immediately to
// both login and logout.
function AuthGate(props) {
  const [authed, setAuthed] = createSignal(pb.authStore.isValid);
  const unsubscribe = pb.authStore.onChange(() => setAuthed(pb.authStore.isValid));
  onCleanup(unsubscribe);

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
        <Route path="/" component={Home} />
        <Route path="/settings" component={Settings} />
      </Router>
    </AuthGate>
  ),
  document.getElementById("app"),
);

