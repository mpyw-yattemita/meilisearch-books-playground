const pattern = /<mark>(.*?)<\/mark>/gs;

/**
 * Renders Meilisearch's highlighted text without dangerouslySetInnerHTML.
 * Only the <mark> tags we asked for become elements; everything else stays text.
 */
export function Highlight({ text }: { text: string }) {
  const parts: React.ReactNode[] = [];
  let last = 0;
  for (const m of text.matchAll(pattern)) {
    parts.push(text.slice(last, m.index));
    parts.push(<mark key={m.index}>{m[1]}</mark>);
    last = m.index + m[0].length;
  }
  parts.push(text.slice(last));
  return <>{parts}</>;
}
