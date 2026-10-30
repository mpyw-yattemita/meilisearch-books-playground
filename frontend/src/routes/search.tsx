import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { getRouteApi, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { type FacetName, type SearchParams, searchBooks, type SortOrder, sortOrders } from "../api";
import { Highlight } from "../Highlight";

const route = getRouteApi("/");

// The URL is the source of truth for the search state, so results are shareable and survive reloads.
export function validateSearch(raw: Record<string, unknown>): SearchParams {
  const str = (v: unknown) => (typeof v === "string" && v !== "" ? v : undefined);
  const num = (v: unknown) => {
    const n = typeof v === "number" ? v : typeof v === "string" && v !== "" ? Number(v) : NaN;
    return Number.isFinite(n) ? n : undefined;
  };
  const list = (v: unknown) => {
    const a = (Array.isArray(v) ? v : [v]).filter((x): x is string => typeof x === "string" && x !== "");
    return a.length > 0 ? a : undefined;
  };
  const sort = str(raw.sort);
  const page = num(raw.page);
  return {
    q: str(raw.q),
    genre: list(raw.genre),
    publisher: list(raw.publisher),
    format: list(raw.format),
    tag: list(raw.tag),
    yearMin: num(raw.yearMin),
    yearMax: num(raw.yearMax),
    priceMin: num(raw.priceMin),
    priceMax: num(raw.priceMax),
    ratingMin: num(raw.ratingMin),
    sort: sort && sort in sortOrders ? (sort as SortOrder) : undefined,
    page: page && page > 1 ? Math.floor(page) : undefined,
  };
}

// The facet name in Meilisearch and the query parameter name differ only for tags.
const facets: { name: FacetName; param: "genre" | "publisher" | "format" | "tag"; label: string }[] = [
  { name: "genre", param: "genre", label: "ジャンル" },
  { name: "format", param: "format", label: "判型" },
  { name: "publisher", param: "publisher", label: "出版社" },
  { name: "tags", param: "tag", label: "タグ" },
];

export function SearchPage() {
  const search = route.useSearch();
  const navigate = route.useNavigate();
  const update = (patch: Partial<SearchParams>) =>
    navigate({ search: (prev) => ({ ...prev, page: undefined, ...patch }), replace: true });

  const { data, error, isFetching } = useQuery({
    queryKey: ["books", search],
    queryFn: ({ signal }) => searchBooks(search, signal),
    placeholderData: keepPreviousData,
  });

  return (
    <div className="search-layout">
      <aside className="filters">
        {facets.map((f) => (
          <FacetList
            key={f.name}
            label={f.label}
            counts={data?.facets[f.name]}
            selected={search[f.param] ?? []}
            onChange={(values) => update({ [f.param]: values.length > 0 ? values : undefined })}
          />
        ))}
        <RangeFilter
          label="刊行年"
          min={search.yearMin}
          max={search.yearMax}
          placeholder={data?.facetStats.publishedYear}
          onChange={(yearMin, yearMax) => update({ yearMin, yearMax })}
        />
        <RangeFilter
          label="価格（円）"
          min={search.priceMin}
          max={search.priceMax}
          placeholder={data?.facetStats.price}
          onChange={(priceMin, priceMax) => update({ priceMin, priceMax })}
        />
        <fieldset>
          <legend>評価</legend>
          <select
            value={search.ratingMin ?? ""}
            onChange={(e) => update({ ratingMin: e.target.value ? Number(e.target.value) : undefined })}
          >
            <option value="">指定なし</option>
            {[4.5, 4, 3.5, 3].map((r) => (
              <option key={r} value={r}>
                ★ {r} 以上
              </option>
            ))}
          </select>
        </fieldset>
        <Link to="/" search={{}} className="clear-filters">
          条件をすべてクリア
        </Link>
      </aside>

      <section className="results">
        <SearchBox initial={search.q ?? ""} onSearch={(q) => update({ q: q || undefined })} />
        <div className="results-bar">
          <span aria-live="polite">
            {data && (
              <>
                <strong>{data.totalHits.toLocaleString()}</strong> 件 / {data.processingTimeMs} ms
              </>
            )}
            {isFetching && <span className="spinner" aria-label="読み込み中" />}
          </span>
          <select value={search.sort ?? "relevance"} onChange={(e) => update({ sort: e.target.value as SortOrder })}>
            {Object.entries(sortOrders).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>

        {error && <p className="error">検索に失敗しました: {error.message}</p>}

        <ol className="hits" style={{ opacity: isFetching ? 0.6 : 1 }}>
          {data?.hits.map((hit) => (
            <li key={hit.id} className="hit">
              <Link to="/books/$id" params={{ id: hit.id }} className="hit-title">
                <Highlight text={hit._formatted.title} />
              </Link>
              <div className="hit-meta">
                <Highlight text={hit._formatted.author} /> ・ {hit.publisher} ・ {hit.publishedYear} 年 ・ {hit.format}
              </div>
              <p className="hit-description">
                <Highlight text={hit._formatted.description} />
              </p>
              <div className="hit-footer">
                <span className="tag">{hit.genre}</span>
                <span>¥{hit.price.toLocaleString()}</span>
                <span>★ {hit.rating.toFixed(1)}</span>
                {search.q && <span className="score">score {hit._rankingScore.toFixed(3)}</span>}
              </div>
            </li>
          ))}
        </ol>
        {data?.hits.length === 0 && <p>見つかりませんでした。</p>}

        {data && data.totalPages > 1 && (
          <Pagination page={data.page} totalPages={data.totalPages} onChange={(page) => update({ page })} />
        )}
      </section>
    </div>
  );
}

function SearchBox({ initial, onSearch }: { initial: string; onSearch: (q: string) => void }) {
  const [value, setValue] = useState(initial);
  // Follow the URL when it changes from outside (back button, links).
  useEffect(() => setValue(initial), [initial]);
  // Search as you type, but not on every keystroke.
  useEffect(() => {
    if (value === initial) return;
    const timer = setTimeout(() => onSearch(value), 200);
    return () => clearTimeout(timer);
  }, [value]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <input
      type="search"
      className="search-box"
      placeholder="タイトル・著者・説明・ISBN で検索（例: 京都 ミステリー / k8s）"
      value={value}
      onChange={(e) => setValue(e.target.value)}
      autoFocus
    />
  );
}

function FacetList(props: {
  label: string;
  counts: Record<string, number> | undefined;
  selected: string[];
  onChange: (values: string[]) => void;
}) {
  const { label, counts = {}, selected, onChange } = props;
  // Keep selected values visible even when the current query matches none of them.
  const entries = Object.entries(counts);
  for (const s of selected) if (!(s in counts)) entries.push([s, 0]);
  entries.sort((a, b) => b[1] - a[1]);

  return (
    <fieldset>
      <legend>{label}</legend>
      {entries.map(([value, count]) => (
        <label key={value} className="facet-option">
          <input
            type="checkbox"
            checked={selected.includes(value)}
            onChange={(e) => onChange(e.target.checked ? [...selected, value] : selected.filter((s) => s !== value))}
          />
          <span>{value}</span>
          <span className="count">{count.toLocaleString()}</span>
        </label>
      ))}
    </fieldset>
  );
}

function RangeFilter(props: {
  label: string;
  min: number | undefined;
  max: number | undefined;
  placeholder: { min: number; max: number } | undefined;
  onChange: (min: number | undefined, max: number | undefined) => void;
}) {
  const { label, min, max, placeholder, onChange } = props;
  const parse = (v: string) => (v === "" ? undefined : Number(v));
  return (
    <fieldset>
      <legend>{label}</legend>
      <div className="range">
        <input
          type="number"
          value={min ?? ""}
          placeholder={placeholder ? String(placeholder.min) : ""}
          onChange={(e) => onChange(parse(e.target.value), max)}
        />
        〜
        <input
          type="number"
          value={max ?? ""}
          placeholder={placeholder ? String(placeholder.max) : ""}
          onChange={(e) => onChange(min, parse(e.target.value))}
        />
      </div>
    </fieldset>
  );
}

function Pagination({ page, totalPages, onChange }: { page: number; totalPages: number; onChange: (page: number) => void }) {
  const pages = [...new Set([1, page - 2, page - 1, page, page + 1, page + 2, totalPages])]
    .filter((p) => p >= 1 && p <= totalPages)
    .sort((a, b) => a - b);
  return (
    <nav className="pagination">
      <button type="button" disabled={page <= 1} onClick={() => onChange(page - 1)}>
        前へ
      </button>
      {pages.map((p, i) => (
        <span key={p}>
          {i > 0 && p - pages[i - 1]! > 1 && <span className="ellipsis">…</span>}
          <button type="button" aria-current={p === page ? "page" : undefined} onClick={() => onChange(p)}>
            {p}
          </button>
        </span>
      ))}
      <button type="button" disabled={page >= totalPages} onClick={() => onChange(page + 1)}>
        次へ
      </button>
    </nav>
  );
}
