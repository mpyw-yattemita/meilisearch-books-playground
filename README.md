# meilisearch-books-playground

A playground for trying [Meilisearch](https://www.meilisearch.com/) with a Go backend.
It searches 50,000 dummy Japanese books.

| Part | Stack |
| --- | --- |
| Search engine | Meilisearch v1.56 (Docker) |
| Backend | Go 1.27, [meilisearch-go](https://github.com/meilisearch/meilisearch-go) v0.36, `net/http` |
| Frontend | React 19, TanStack Router, TanStack Query, TypeScript 7, Vite 8 |
| Go linter | [declscope](https://github.com/mpyw/declscope) (`qualify: ondemand`, `exported: true`) |

## Quick start

Run each of the last three commands in its own terminal.

```bash
make up      # Meilisearch on :7700
make seed    # 50,000 books
make server  # API on :8080
make web     # UI on http://localhost:5173
```

Measured on an Apple Silicon Mac:

```console
$ make seed
2026/10/30 14:20:34 loading 50000 books (batch size 10000)...
2026/10/30 14:20:36 done in 1.944s
```

> [!TIP]
> `MEILI_ENV=development` turns on the Meilisearch search preview.
> Open http://localhost:7700 and enter the master key from `compose.yaml`.

## What it tries

| Feature | Where | Try |
| --- | --- | --- |
| Japanese tokenization | `localizedAttributes` in [`bookindex.go`](backend/internal/bookindex/bookindex.go) | `京都` matches `京都の館`. `都` alone matches nothing |
| Synonyms | Same file | `k8s`, `golang`, `ts`, `postgres` |
| Highlighting | [`search.go`](backend/internal/bookindex/search.go), [`Highlight.tsx`](frontend/src/Highlight.tsx) | Matches are wrapped in `<mark>` |
| Disjunctive facets | `Search` in `search.go` | Check one genre. The other genres still show their counts |
| Filters | [`filter.go`](backend/internal/bookindex/filter.go) | Year, price and rating ranges |
| Sorting | [`sort.go`](backend/internal/bookindex/sort.go) | Newest, price, rating, review count |
| Page-based pagination | `page` and `hitsPerPage` | `maxTotalHits` is raised to 100,000 so every page is reachable |
| Ranking score | `showRankingScore` | Shown on each hit while a query is entered |

> [!NOTE]
> Disjunctive facets cost one extra query per facet that has a selection.
> All queries go in one multi-search request.

## API

| Endpoint | Description |
| --- | --- |
| `GET /api/books` | Search. Parameters are below |
| `GET /api/books/{id}` | One book. `404` if it does not exist |
| `GET /api/health` | Meilisearch health |

| Parameter | Meaning |
| --- | --- |
| `q` | Full-text query |
| `genre`, `publisher`, `format`, `tag` | Facet values. Repeat a parameter to select several |
| `yearMin`, `yearMax`, `priceMin`, `priceMax`, `ratingMin` | Range filters |
| `sort` | `relevance` (default), `newest`, `oldest`, `price_asc`, `price_desc`, `rating`, `popular` |
| `page`, `hitsPerPage` | Page starts at 1. `hitsPerPage` is at most 100 |

Values inside one facet are OR-ed. Different facets are AND-ed.

```bash
curl -sG localhost:8080/api/books --data-urlencode q=京都 --data-urlencode genre=ミステリー
```

## Layout

| Path | Contents |
| --- | --- |
| `compose.yaml` | Meilisearch |
| `backend/cmd/seed` | Generates books and loads them. `-dump` prints them as JSON instead |
| `backend/cmd/server` | The HTTP API |
| `backend/internal/bookgen` | Deterministic dummy data. The same `-seed` gives the same books |
| `backend/internal/bookindex` | Index settings, loading, and search |
| `frontend/src/routes` | Search page and book page |

The frontend dev server proxies `/api` to `localhost:8080`.

## Configuration

| Variable | Default |
| --- | --- |
| `MEILI_HOST` | `http://localhost:7700` |
| `MEILI_MASTER_KEY` | `local-dev-master-key-change-me` |

> [!WARNING]
> The API server uses the master key for simplicity.
> A real deployment should give it a search-only API key.

## Lint

```bash
make lint
```

It runs `go vet`, `declscope shrink`, `declscope` and `tsc`.
