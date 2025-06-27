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
import React, {
  PropsWithChildren,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

type PaginationContextType<O extends DescMessage> = UseInfiniteQueryResult<
  InfiniteData<MessageShape<O>>,
  ConnectError
> & {
  page: MessageShape<O>;
  consoleFetchNextPage: () => void;
  consoleFetchPreviousPage: () => void;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const PaginationContext = createContext<PaginationContextType<any> | null>(
  null,
);

export function PaginationProvider<O extends DescMessage>({
  children,
  query,
}: {
  query: UseInfiniteQueryResult<InfiniteData<MessageShape<O>>, ConnectError> & {
    page: MessageShape<O>;
    consoleFetchNextPage: () => void;
    consoleFetchPreviousPage: () => void;
  };
} & PropsWithChildren) {
  return (
    <PaginationContext.Provider value={query}>
      {children}
    </PaginationContext.Provider>
  );
}

export function usePagination() {
  const context = useContext(PaginationContext);
  if (!context)
    throw new Error(
      "usePaginationContext must be used inside a PaginationProvider",
    );

  return context;
}

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

  const {
    data,
    hasNextPage: queryHasNextPage,
    fetchNextPage,
    isFetching,
    isSuccess,
    ...result
  } = useInfiniteQuery(schema, input, {
    ...queryOptions,
    retry: queryOptions.retry ?? 3,
  });

  function consoleFetchNextPage() {
    const nextPage = currentPage + 1;

    const prefetchedPage = data?.pages[nextPage - 1];

    if (prefetchedPage) {
      setCurrentPage(nextPage);
      setPage(prefetchedPage);
    } else if (queryHasNextPage) {
      fetchNextPage().then(() => {
        setCurrentPage(nextPage);
      });
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

    const prefetchedPage = data.pages[previousPage - 1];
    if (prefetchedPage) {
      setCurrentPage(previousPage);
      setPage(prefetchedPage);
    }
  }

  useEffect(() => {
    if (data?.pages.length && currentPage === 0) {
      setCurrentPage(1);
      setPage(data.pages[0]);
    }
  }, [currentPage, data, setCurrentPage, setPage]);

  useEffect(() => {
    if (data && isSuccess && !isFetching) {
      const newPageCount = data.pages.length;

      setHasNextPage(currentPage < newPageCount || queryHasNextPage);

      const pageResults = data.pages[currentPage - 1];
      if (pageResults) {
        setPage(pageResults);
      }
    }
  }, [
    currentPage,
    data,
    isSuccess,
    isFetching,
    queryHasNextPage,
    setHasNextPage,
    setPage,
  ]);

  useEffect(() => {
    if (!data) {
      return;
    }

    setHasPreviousPage(currentPage > 1);
    setHasNextPage(currentPage < data.pages.length || queryHasNextPage);
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
