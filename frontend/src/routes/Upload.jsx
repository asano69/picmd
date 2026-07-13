import { createSignal, Show } from "solid-js";
import NavBar from "../components/NavBar";
import Button from "../components/Button";
import pb from "../lib/pb";

// formatSize renders a byte count as a short human-readable string.
function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

// Upload lets the user paste, drop, or select an image, sends it to the
// "images" collection (compression happens server-side, see
// internal/hooks/images.go), and shows a Markdown snippet for the result.
export default function Upload() {
  const [file, setFile] = createSignal(null);
  const [previewUrl, setPreviewUrl] = createSignal("");
  const [status, setStatus] = createSignal("");
  const [uploading, setUploading] = createSignal(false);
  const [result, setResult] = createSignal(null);
  const [copied, setCopied] = createSignal(false);

  let fileInputRef;

  const pickFile = (f) => {
    if (!f || !f.type.startsWith("image/")) {
      setStatus("Only image files are supported.");
      return;
    }
    setFile(f);
    setPreviewUrl(URL.createObjectURL(f));
    setResult(null);
    setStatus("");
  };

  const clear = () => {
    setFile(null);
    setPreviewUrl("");
    setResult(null);
    setStatus("");
  };

  const upload = async () => {
    const f = file();
    if (!f) return;
    setUploading(true);
    setStatus("Uploading…");
    try {
      const formData = new FormData();
      formData.append("image", f, f.name || `clipboard-${Date.now()}.png`);
      const record = await pb.collection("images").create(formData);
      setResult({
        url: pb.files.getURL(record, record.image),
        filename: record.filename,
        filesize: record.filesize,
      });
      setCopied(false);
      setStatus("Upload complete.");
    } catch (err) {
      setStatus(`Upload failed: ${err.message}`);
    } finally {
      setUploading(false);
    }
  };

  const copyMarkdown = () => {
    const r = result();
    if (!r) return;
    navigator.clipboard.writeText(`![](${r.url})`);
    setCopied(true);
  };

  // Pasting anywhere on the page picks up a clipboard image, matching the
  // old upload.html behavior.
  const onPaste = (e) => {
    const items = e.clipboardData?.items;
    if (!items) return;
    for (const item of items) {
      if (item.kind === "file" && item.type.startsWith("image/")) {
        pickFile(item.getAsFile());
        break;
      }
    }
  };

  const onDrop = (e) => {
    e.preventDefault();
    const f = e.dataTransfer.files[0];
    if (f) pickFile(f);
  };

  return (
    <div
      class="mx-auto flex min-h-screen w-full max-w-xl flex-col items-center bg-[var(--color-bg)] px-6 py-12 text-[var(--color-text)]"
      onPaste={onPaste}
    >
      <NavBar />

      <Show
        when={file()}
        fallback={
          <div
            tabIndex="0"
            role="button"
            aria-label="Paste, drop, or click to select an image"
            class="flex w-full cursor-pointer flex-col items-center rounded-md border border-dashed border-[var(--color-border-soft)] bg-[var(--color-panel)] px-8 py-14 text-center"
            onDragOver={(e) => e.preventDefault()}
            onDrop={onDrop}
            onClick={() => fileInputRef.click()}
            onKeyDown={(e) => (e.key === "Enter" || e.key === " ") && fileInputRef.click()}
          >
            <p class="text-sm leading-loose">
              <span class="font-semibold">Ctrl+V</span> to paste from clipboard
              <br />
              or drag &amp; drop an image here
              <br />
              or <span class="font-semibold">click</span> to select a file
            </p>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*,.svg"
              class="hidden"
              onChange={(e) => {
                if (e.target.files[0]) pickFile(e.target.files[0]);
                e.target.value = ""; // allow re-selecting the same file
              }}
            />
          </div>
        }
      >
        <div class="w-full">
          <img
            src={previewUrl()}
            alt="Preview"
            class="w-full rounded-md border border-[var(--color-border-soft)]"
          />
          <p class="mt-2 text-sm opacity-70">
            {file().name || "clipboard-image"} · {file().type} · {formatSize(file().size)}
          </p>

          <Show when={!result()}>
            <div class="mt-3 flex flex-wrap">
              <Button value={uploading() ? "Uploading…" : "Upload"} disabled={uploading()} onClick={upload} />
              <Button value="Clear" disabled={uploading()} onClick={clear} />
            </div>
          </Show>

          <Show when={result()}>
            {(r) => (
              <div class="mt-3 rounded-md border border-[var(--color-border-soft)] bg-[var(--color-panel)] p-4 text-sm leading-loose">
                <div class="flex justify-between border-b border-[var(--color-border-soft)] pb-1">
                  <span class="opacity-70">Filename</span>
                  <span class="font-medium">{r().filename}</span>
                </div>
                <div class="flex justify-between py-1">
                  <span class="opacity-70">Size</span>
                  <span class="font-medium">{formatSize(r().filesize)}</span>
                </div>
                <div class="mt-2 flex flex-wrap">
                  <Button value={copied() ? "Copied!" : "Copy Markdown"} onClick={copyMarkdown} />
                  <Button value="Upload Another" onClick={clear} />
                </div>
              </div>
            )}
          </Show>
        </div>
      </Show>

      <Show when={status()}>
        <p class="mt-4 text-sm opacity-70">{status()}</p>
      </Show>
    </div>
  );
}
