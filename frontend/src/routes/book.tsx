import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { getRouteApi, Link, useCanGoBack, useRouter } from "@tanstack/react-router";
import { getBook } from "../api";

export const bookQueryOptions = (id: number) =>
  queryOptions({
    queryKey: ["book", id],
    queryFn: ({ signal }) => getBook(id, signal),
  });

const route = getRouteApi("/books/$id");

export function BookPage() {
  const { id } = route.useParams();
  const { data: book } = useSuspenseQuery(bookQueryOptions(id));
  const router = useRouter();
  const canGoBack = useCanGoBack();

  return (
    <article className="book-detail">
      {canGoBack ? (
        <button type="button" className="link-button" onClick={() => router.history.back()}>
          ← 検索結果に戻る
        </button>
      ) : (
        <Link to="/">← 検索へ</Link>
      )}
      <h1>{book.title}</h1>
      <p className="book-author">{book.author}</p>
      <p>{book.description}</p>
      <dl>
        <dt>ジャンル</dt>
        <dd>
          <Link to="/" search={{ genre: [book.genre] }}>
            {book.genre}
          </Link>
        </dd>
        <dt>出版社</dt>
        <dd>
          <Link to="/" search={{ publisher: [book.publisher] }}>
            {book.publisher}
          </Link>
        </dd>
        <dt>判型</dt>
        <dd>{book.format}</dd>
        <dt>タグ</dt>
        <dd className="tags">
          {book.tags.map((t) => (
            <Link key={t} to="/" search={{ tag: [t] }} className="tag">
              {t}
            </Link>
          ))}
        </dd>
        <dt>刊行年</dt>
        <dd>{book.publishedYear} 年</dd>
        <dt>ページ数</dt>
        <dd>{book.pages} ページ</dd>
        <dt>価格</dt>
        <dd>¥{book.price.toLocaleString()}</dd>
        <dt>評価</dt>
        <dd>
          ★ {book.rating.toFixed(1)}（{book.reviewCount.toLocaleString()} 件）
        </dd>
        <dt>ISBN</dt>
        <dd>{book.isbn}</dd>
      </dl>
    </article>
  );
}
