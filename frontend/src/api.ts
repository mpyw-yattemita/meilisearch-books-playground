// Types and fetchers for the Go API (backend/cmd/server).

export type Book = {
  id: number;
  isbn: string;
  title: string;
  author: string;
  publisher: string;
  genre: string;
  format: string;
  tags: string[];
  description: string;
  publishedYear: number;
  pages: number;
  price: number;
  rating: number;
  reviewCount: number;
};

export type SearchHit = Book & {
  /** Copies of searchable fields with matches wrapped in <mark>...</mark>. */
  _formatted: Pick<Book, "title" | "author" | "description">;
  _rankingScore: number;
};

export type FacetName = "genre" | "publisher" | "format" | "tags";

export type SearchResult = {
  hits: SearchHit[];
  query: string;
  page: number;
  hitsPerPage: number;
  totalHits: number;
  totalPages: number;
  processingTimeMs: number;
  facets: Partial<Record<FacetName, Record<string, number>>>;
  facetStats: Partial<Record<"publishedYear" | "price", { min: number; max: number }>>;
};

export const sortOrders = {
  relevance: "関連度順",
  newest: "新しい順",
  oldest: "古い順",
  price_asc: "価格の安い順",
  price_desc: "価格の高い順",
  rating: "評価の高い順",
  popular: "レビュー数順",
} as const;

export type SortOrder = keyof typeof sortOrders;

export type SearchParams = {
  q?: string;
  genre?: string[];
  publisher?: string[];
  format?: string[];
  tag?: string[];
  yearMin?: number;
  yearMax?: number;
  priceMin?: number;
  priceMax?: number;
  ratingMin?: number;
  sort?: SortOrder;
  page?: number;
};

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message);
  }
}

async function getJSON<T>(path: string, signal?: AbortSignal): Promise<T> {
  const res = await fetch(path, { signal });
  const body: unknown = await res.json().catch(() => null);
  if (!res.ok) {
    const message =
      body && typeof body === "object" && "error" in body && typeof body.error === "string"
        ? body.error
        : res.statusText;
    throw new ApiError(res.status, message);
  }
  return body as T;
}

export function searchBooks(params: SearchParams, signal?: AbortSignal): Promise<SearchResult> {
  const qs = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === "") continue;
    for (const v of Array.isArray(value) ? value : [value]) qs.append(key, String(v));
  }
  return getJSON(`/api/books?${qs}`, signal);
}

export function getBook(id: number, signal?: AbortSignal): Promise<Book> {
  return getJSON(`/api/books/${id}`, signal);
}
