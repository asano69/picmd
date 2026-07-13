import { createSignal, For, Show } from "solid-js";
import NavBar from "../components/NavBar";
import Button from "../components/Button";
import pb from "../lib/pb";

// formatSize renders a byte count as a short human-readable string.
function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

let nextId = 0;

// Upload lets the user paste, drop, or select one or more images. Each
// image becomes its own "images" record (compression happens
// server-side per record, see internal/hooks/images.go), so uploading
// several files just means running several independent uploads here —
// no backend change is needed to support multi-file selection.
export default function Upload() {
  const [items, setItems] = createSignal([]);
  const [status, setStatus] = createSignal("");
  const [allCopied, setAllCopied] = createSignal(false);

  let fileInputRef;

  const addFiles = (fileList) => {
    const files = [...fileList].filter((f) => f.type.startsWith("image/"));
    if (files.length === 0) {
      setStatus("Only image files are supported.");
      return;
    }
    const newItems = files.map((file) => ({
      id: nextId++,
      file,
      previewUrl: URL.createObjectURL(file),
      status: "pending", // pending | uploading | done | error
      result: null,
      error: "",
    }));
    setItems((prev) => [...prev, ...newItems]);
    setStatus("");
    setAllCopied(false);
  };

  const updateItem = (id, patch) =>
    setItems((prev) => prev.map((it) => (it.id === id ? { ...it, ...patch } : it)));

  const removeItem = (id) => {
    const it = items().find((i) => i.id === id);
    if (it) URL.revokeObjectURL(it.previewUrl);
    setItems((prev) => prev.filter((i) => i.id !== id));
  };

  const uploadOne = async (item) => {
    updateItem(item.id, { status: "uploading", error: "" });
    try {
      const formData = new FormData();
      formData.append(
        "image",
        item.file,
        item.file.name || `clipboard-${item.id}.png`,
      );
      // Each upload gets its own requestKey so PocketBase's SDK-level
      // auto-cancellation (which assumes concurrent calls to the same
      // endpoint are duplicates) doesn't abort sibling uploads.
      const record = await pb.collection("images").create(formData, {
        requestKey: `upload-${item.id}`,
      });
      updateItem(item.id, {
        status: "done",
        result: {
          url: new URL(`/i/${record.uuid}`, window.location.origin).href,
          filename: record.filename,
          filesize: record.filesize,
        },
      });
    } catch (err) {
      updateItem(item.id, { status: "error", error: err.message });
    }
  };

  const uploadAll = () => {
    const pending = items().filter(
      (it) => it.status === "pending" || it.status === "error",
    );
    if (pending.length === 0) return;
    setStatus(`Uploading ${pending.length} image${pending.length > 1 ? "s" : ""}…`);
    Promise.all(pending.map(uploadOne)).then(() => setStatus(""));
  };

  const clearAll = () => {
    items().forEach((it) => URL.revokeObjectURL(it.previewUrl));
    setItems([]);
    setStatus("");
    setAllCopied(false);
  };

  const copyAllMarkdown = () => {
    const done = items().filter((it) => it.status === "done");
    if (done.length === 0) return;
    const markdown = done.map((it) => `![](${it.result.url})`).join("\n");
    navigator.clipboard.writeText(markdown);
    setAllCopied(true);
  };

  // Pasting anywhere on the page picks up every image on the clipboard,
  // instead of stopping at the first one like the single-image version did.
  const onPaste = (e) => {
    const clipboardItems = e.clipboardData?.items;
    if (!clipboardItems) return;
    const files = [...clipboardItems]
      .filter((it) => it.kind === "file" && it.type.startsWith("image/"))
      .map((it) => it.getAsFile());
    if (files.length > 0) addFiles(files);
  };

  const onDrop = (e) => {
    e.preventDefault();
    addFiles(e.dataTransfer.files);
  };

  const hasItems = () => items().length > 0;
  const hasPending = () =>
    items().some((it) => it.status === "pending" || it.status === "error");
  const hasDone = () => items().some((it) => it.status === "done");

  return (
    <div
      class="mx-auto flex min-h-screen w-full max-w-xl flex-col items-center bg-[var(--color-bg)] px-6 py-12 text-[var(--color-text)]"
      onPaste={onPaste}
    >
      <NavBar />

      {/* Kept outside the Show branches so "Add More" can reuse it too. */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*,.svg"
        multiple
        class="hidden"
        onChange={(e) => {
          if (e.target.files.length) addFiles(e.target.files);
          e.target.value = ""; // allow re-selecting the same file(s)
        }}
      />

      <Show
        when={hasItems()}
        fallback={
          <div
            tabIndex="0"
            role="button"
            aria-label="Paste, drop, or click to select images"
            class="flex w-full cursor-pointer flex-col items-center rounded-md border border-dashed border-[var(--color-border-soft)] bg-[var(--color-panel)] px-8 py-14 text-center"
            onDragOver={(e) => e.preventDefault()}
            onDrop={onDrop}
            onClick={() => fileInputRef.click()}
            onKeyDown={(e) =>
              (e.key === "Enter" || e.key === " ") && fileInputRef.click()
            }
          >
            <p class="text-sm leading-loose">
              <span class="font-semibold">Ctrl+V</span> to paste from clipboard
              <br />
              or drag &amp; drop images here
              <br />
              or <span class="font-semibold">click</span> to select files
            </p>
          </div>
        }
      >
        <div class="w-full">
          <div class="flex flex-col gap-3">
            <For each={items()}>
              {(item) => (
                <div class="flex items-center gap-3 rounded-md border border-[var(--color-border-soft)] bg-[var(--color-panel)] p-3">
                  <img
                    src={item.previewUrl}
                    alt="Preview"
                    class="h-16 w-16 flex-none rounded object-cover"
                  />
                  <div class="min-w-0 flex-1 text-sm">
                    <p class="truncate">{item.file.name || "clipboard-image"}</p>
                    <p class="opacity-70">
                      {formatSize(item.file.size)}
                      {item.status === "uploading" && " · Uploading…"}
                      {item.status === "done" &&
                        ` · ${formatSize(item.result.filesize)} compressed`}
                      {item.status === "error" && ` · Failed: ${item.error}`}
                    </p>
                  </div>
                  <Show when={item.status === "done"}>
                    <Button
                      value="Copy"
                      onClick={() =>
                        navigator.clipboard.writeText(`![](${item.result.url})`)
                      }
                    />
                  </Show>
                  <Show when={item.status !== "uploading"}>
                    <Button value="Remove" onClick={() => removeItem(item.id)} />
                  </Show>
                </div>
              )}
            </For>
          </div>

          <div class="mt-3 flex flex-wrap">
            <Show when={hasPending()}>
              <Button value="Upload All" onClick={uploadAll} />
            </Show>
            <Show when={hasDone()}>
              <Button
                value={allCopied() ? "Copied!" : "Copy All as Markdown"}
                onClick={copyAllMarkdown}
              />
            </Show>
            <Button value="Add More" onClick={() => fileInputRef.click()} />
            <Button value="Clear All" onClick={clearAll} />
          </div>
        </div>
      </Show>

      <Show when={status()}>
        <p class="mt-4 text-sm opacity-70">{status()}</p>
      </Show>
    </div>
  );
}
