import {
  DescMessage,
  DescMethodUnary,
  MessageInitShape,
  MessageShape,
} from "@bufbuild/protobuf";
import { ConnectError } from "@connectrpc/connect";
import {
  UseInfiniteQueryOptions,
  useInfiniteQuery,
} from "@connectrpc/connect-query";
import {
  InfiniteData,
  SkipToken,
  UseInfiniteQueryResult,
} from "@tanstack/react-query";
import { useEffect, useState } from "react";

export function usePaginatedInfiniteQuery<
  I extends DescMessage,
  O extends DescMessage,
  ParamKey extends keyof MessageInitShape<I>,
>(
  schema: DescMethodUnary<I, O>,
  input:
    | SkipToken
    | (MessageInitShape<I> & Required<Pick<MessageInitShape<I>, ParamKey>>),
  queryOptions: UseInfiniteQueryOptions<I, O, ParamKey>,
): UseInfiniteQueryResult<InfiniteData<MessageShape<O>>, ConnectError> & {
  page: MessageShape<O>;
  consoleFetchNextPage: () => void;
  consoleFetchPreviousPage: () => void;
} {
  const [currentPage, setCurrentPage] = useState(0);
  const [hasNextPage, setHasNextPage] = useState(false);
  const [hasPreviousPage, setHasPreviousPage] = useState(false);
  const [page, setPage] = useState<MessageShape<O>>({} as MessageShape<O>);
  const [fetchedPageCount, setFetchedPageCount] = useState(0);

  const {
    data,
    hasNextPage: queryHasNextPage,
    fetchNextPage,
    isFetching,
    isSuccess,
    ...result
  } = useInfiniteQuery(schema, input, queryOptions);

  function consoleFetchNextPage() {
    if (!data || (currentPage >= fetchedPageCount && !queryHasNextPage)) {
      return;
    }

    const nextPage = currentPage + 1;
    if (nextPage > fetchedPageCount) {
      fetchNextPage();
    } else if (data.pages[nextPage - 1]) {
      setCurrentPage(nextPage);
    }
  }

  function consoleFetchPreviousPage() {
    if (!data || currentPage < 2) {
      return;
    }

    const previousPage = currentPage - 1;
    if (previousPage < 1) {
      return;
    }

    if (data.pages[previousPage - 1]) {
      setCurrentPage(previousPage);
    }
  }

  useEffect(() => {
    if (data && isSuccess && !isFetching) {
      const newPageCount = data.pages.length;

      setHasNextPage(newPageCount > currentPage || queryHasNextPage);

      // Only run if fetching next page
      if (newPageCount > fetchedPageCount) {
        setFetchedPageCount(data.pages.length);

        if (newPageCount > currentPage) {
          setCurrentPage(newPageCount);
        }
      }
    }
  }, [
    currentPage,
    data,
    isSuccess,
    isFetching,
    fetchedPageCount,
    queryHasNextPage,
    setHasNextPage,
    setFetchedPageCount,
    setCurrentPage,
  ]);

  useEffect(() => {
    if (!data) {
      return;
    }

    setHasPreviousPage(currentPage > 1);
    setHasNextPage(currentPage < data.pages.length || queryHasNextPage);

    const pageResults = data.pages[currentPage - 1];
    if (!pageResults) {
      return;
    }
    setPage(pageResults);
  }, [currentPage, data, queryHasNextPage, setPage, setHasPreviousPage]);

  return {
    ...result,
    consoleFetchNextPage,
    consoleFetchPreviousPage,
    data,
    hasNextPage,
    hasPreviousPage,
    isFetching,
    isSuccess,
    page,
  } as UseInfiniteQueryResult<
    InfiniteData<MessageShape<O>, unknown>,
    ConnectError
  > & {
    page: MessageShape<O>;
    consoleFetchNextPage: () => void;
    consoleFetchPreviousPage: () => void;
  };
}
