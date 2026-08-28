import {
  useQuery,
  useMutation,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { boardsApi } from "../api/boardsApi";
import type {
  BoardView,
  CreateBoardRequest,
  UpdateBoardRequest,
  PagedModelBoardView,
  BoardColumnView,
  CreateBoardColumnRequest,
  UpdateBoardColumnRequest,
  ReorderBoardColumnsRequest,
  CardSummaryView,
  CardView,
  CreateCardRequest,
  UpdateCardRequest,
  MoveCardRequest,
} from "~types/api";

// ── Query Keys ───────────────────────────────────────────────
export const boardKeys = {
  all: ["boards"] as const,
  lists: (projectId: string) => [...boardKeys.all, projectId] as const,
  defaultBoard: (projectId: string) => [...boardKeys.lists(projectId), "default"] as const,
  detail: (projectId: string, boardId: string) =>
    [...boardKeys.lists(projectId), boardId] as const,
  columns: (projectId: string, boardId: string) =>
    [...boardKeys.detail(projectId, boardId), "columns"] as const,
  cards: (projectId: string, boardId: string) =>
    [...boardKeys.detail(projectId, boardId), "cards"] as const,
  cardDetail: (projectId: string, boardId: string, cardId: string) =>
    [...boardKeys.cards(projectId, boardId), cardId] as const,
};

// ── Board Hooks ──────────────────────────────────────────────
export const useBoardsQuery = (
  projectId: string,
  page = 0,
  size = 100
): UseQueryResult<PagedModelBoardView, Error> => {
  return useQuery({
    queryKey: [...boardKeys.lists(projectId), { page, size }],
    queryFn: () => boardsApi.getAllBoards(projectId, page, size),
    enabled: Boolean(projectId),
  });
};

export const useDefaultBoardQuery = (
  projectId: string
): UseQueryResult<BoardView, Error> => {
  return useQuery({
    queryKey: boardKeys.defaultBoard(projectId),
    queryFn: () => boardsApi.getDefaultBoard(projectId),
    enabled: Boolean(projectId),
  });
};

export const useBoardQuery = (
  projectId: string,
  boardId?: string
): UseQueryResult<BoardView, Error> => {
  return useQuery({
    queryKey: boardKeys.detail(projectId, boardId || ""),
    queryFn: () => boardsApi.getBoardById(projectId, boardId!),
    enabled: Boolean(projectId && boardId),
  });
};

export const useCreateBoardMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateBoardRequest) =>
      boardsApi.createBoard(projectId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: boardKeys.lists(projectId) });
    },
  });
};

export const useUpdateBoardMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      boardId,
      payload,
    }: {
      boardId: string;
      payload: UpdateBoardRequest;
    }) => boardsApi.updateBoard(projectId, boardId, payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: boardKeys.lists(projectId) });
      queryClient.invalidateQueries({
        queryKey: boardKeys.detail(projectId, variables.boardId),
      });
    },
  });
};

export const useDeleteBoardMutation = (projectId: string) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (boardId: string) => boardsApi.deleteBoard(projectId, boardId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: boardKeys.lists(projectId) });
    },
  });
};

// ── Column Hooks ─────────────────────────────────────────────
export const useBoardColumnsQuery = (
  projectId: string,
  boardId?: string
): UseQueryResult<BoardColumnView[], Error> => {
  return useQuery({
    queryKey: boardId ? boardKeys.columns(projectId, boardId) : [],
    queryFn: () => boardsApi.getColumns(projectId, boardId!),
    enabled: Boolean(projectId && boardId),
  });
};

export const useCreateColumnMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateBoardColumnRequest) =>
      boardsApi.createColumn(projectId, boardId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.columns(projectId, boardId),
      });
    },
  });
};

export const useUpdateColumnMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      columnId,
      payload,
    }: {
      columnId: string;
      payload: UpdateBoardColumnRequest;
    }) => boardsApi.updateColumn(projectId, boardId, columnId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.columns(projectId, boardId),
      });
    },
  });
};

export const useDeleteColumnMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (columnId: string) =>
      boardsApi.deleteColumn(projectId, boardId, columnId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.columns(projectId, boardId),
      });
      queryClient.invalidateQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
    },
  });
};

export const useReorderColumnsMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReorderBoardColumnsRequest) =>
      boardsApi.reorderColumns(projectId, boardId, payload),
    onMutate: async (newOrder) => {
      await queryClient.cancelQueries({
        queryKey: boardKeys.columns(projectId, boardId),
      });
      const previousColumns = queryClient.getQueryData<BoardColumnView[]>(
        boardKeys.columns(projectId, boardId)
      );
      if (previousColumns) {
        const orderMap = new Map(
          newOrder.positions.map((id, index) => [id, index])
        );
        const reordered = [...previousColumns].sort((a, b) => {
          const posA = orderMap.get(a.id) ?? a.position;
          const posB = orderMap.get(b.id) ?? b.position;
          return posA - posB;
        });
        queryClient.setQueryData(
          boardKeys.columns(projectId, boardId),
          reordered
        );
      }
      return { previousColumns };
    },
    onError: (_err, _newOrder, context) => {
      if (context?.previousColumns) {
        queryClient.setQueryData(
          boardKeys.columns(projectId, boardId),
          context.previousColumns
        );
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.columns(projectId, boardId),
      });
    },
  });
};

// ── Card Hooks ───────────────────────────────────────────────
export const useBoardCardsQuery = (
  projectId: string,
  boardId?: string
): UseQueryResult<CardSummaryView[], Error> => {
  return useQuery({
    queryKey: boardKeys.cards(projectId, boardId || ""),
    queryFn: () => boardsApi.getCards(projectId, boardId!),
    enabled: Boolean(projectId && boardId),
  });
};

export const useCardDetailQuery = (
  projectId: string,
  boardId?: string,
  cardId?: string
): UseQueryResult<CardView, Error> => {
  return useQuery({
    queryKey: boardKeys.cardDetail(projectId, boardId || "", cardId || ""),
    queryFn: () => boardsApi.getCardById(projectId, boardId!, cardId!),
    enabled: Boolean(projectId && boardId && cardId),
  });
};

export const useCreateCardMutation = (
  projectId: string,
  boardId: string,
  columnId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateCardRequest) =>
      boardsApi.createCard(projectId, boardId, columnId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
    },
  });
};

export const useUpdateCardMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      cardId,
      payload,
    }: {
      cardId: string;
      payload: UpdateCardRequest;
    }) => boardsApi.updateCard(projectId, boardId, cardId, payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
      queryClient.invalidateQueries({
        queryKey: boardKeys.cardDetail(
          projectId,
          boardId,
          variables.cardId
        ),
      });
    },
  });
};

export const useDeleteCardMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (cardId: string) =>
      boardsApi.deleteCard(projectId, boardId, cardId),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
    },
  });
};

export const useMoveCardMutation = (
  projectId: string,
  boardId: string
) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      cardId,
      payload,
    }: {
      cardId: string;
      payload: MoveCardRequest;
    }) => boardsApi.moveCard(projectId, boardId, cardId, payload),
    onMutate: async ({ cardId, payload }) => {
      await queryClient.cancelQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
      const previousCards = queryClient.getQueryData<CardSummaryView[]>(
        boardKeys.cards(projectId, boardId)
      );

      if (previousCards) {
        const updatedCards = previousCards.map((card) => {
          if (card.id === cardId) {
            return {
              ...card,
              columnId: payload.targetColumnId,
              position: payload.position,
            };
          }
          return card;
        });

        queryClient.setQueryData(
          boardKeys.cards(projectId, boardId),
          updatedCards
        );
      }

      return { previousCards };
    },
    onError: (_err, _variables, context) => {
      if (context?.previousCards) {
        queryClient.setQueryData(
          boardKeys.cards(projectId, boardId),
          context.previousCards
        );
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: boardKeys.cards(projectId, boardId),
      });
    },
  });
};
