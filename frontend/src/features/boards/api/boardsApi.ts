import { apiClient } from "@/api/client";
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

export const boardsApi = {
  // ── Boards ──────────────────────────────────────────────────
  getAllBoards: async (
    projectId: string,
    page = 0,
    size = 100
  ): Promise<PagedModelBoardView> => {
    const response = await apiClient.get<PagedModelBoardView>(
      `/projects/${projectId}/boards`,
      { params: { page, size } }
    );
    return response.data;
  },

  getDefaultBoard: async (projectId: string): Promise<BoardView> => {
    const response = await apiClient.get<BoardView>(
      `/projects/${projectId}/boards/default`
    );
    return response.data;
  },

  getBoardById: async (
    projectId: string,
    boardId: string
  ): Promise<BoardView> => {
    const response = await apiClient.get<BoardView>(
      `/projects/${projectId}/boards/${boardId}`
    );
    return response.data;
  },

  createBoard: async (
    projectId: string,
    payload: CreateBoardRequest
  ): Promise<BoardView> => {
    const response = await apiClient.post<BoardView>(
      `/projects/${projectId}/boards`,
      payload
    );
    return response.data;
  },

  updateBoard: async (
    projectId: string,
    boardId: string,
    payload: UpdateBoardRequest
  ): Promise<BoardView> => {
    const response = await apiClient.patch<BoardView>(
      `/projects/${projectId}/boards/${boardId}`,
      payload
    );
    return response.data;
  },

  deleteBoard: async (projectId: string, boardId: string): Promise<void> => {
    await apiClient.delete(`/projects/${projectId}/boards/${boardId}`);
  },

  // ── Columns ─────────────────────────────────────────────────
  getColumns: async (
    projectId: string,
    boardId: string
  ): Promise<BoardColumnView[]> => {
    const response = await apiClient.get<BoardColumnView[]>(
      `/projects/${projectId}/boards/${boardId}/columns`
    );
    return response.data;
  },

  getColumnById: async (
    projectId: string,
    boardId: string,
    columnId: string
  ): Promise<BoardColumnView> => {
    const response = await apiClient.get<BoardColumnView>(
      `/projects/${projectId}/boards/${boardId}/columns/${columnId}`
    );
    return response.data;
  },

  createColumn: async (
    projectId: string,
    boardId: string,
    payload: CreateBoardColumnRequest
  ): Promise<BoardColumnView> => {
    const response = await apiClient.post<BoardColumnView>(
      `/projects/${projectId}/boards/${boardId}/columns`,
      payload
    );
    return response.data;
  },

  updateColumn: async (
    projectId: string,
    boardId: string,
    columnId: string,
    payload: UpdateBoardColumnRequest
  ): Promise<BoardColumnView> => {
    const response = await apiClient.patch<BoardColumnView>(
      `/projects/${projectId}/boards/${boardId}/columns/${columnId}`,
      payload
    );
    return response.data;
  },

  deleteColumn: async (
    projectId: string,
    boardId: string,
    columnId: string
  ): Promise<void> => {
    await apiClient.delete(
      `/projects/${projectId}/boards/${boardId}/columns/${columnId}`
    );
  },

  reorderColumns: async (
    projectId: string,
    boardId: string,
    payload: ReorderBoardColumnsRequest
  ): Promise<BoardColumnView[]> => {
    const response = await apiClient.patch<BoardColumnView[]>(
      `/projects/${projectId}/boards/${boardId}/columns/reorder`,
      payload
    );
    return response.data;
  },

  // ── Cards ───────────────────────────────────────────────────
  getCards: async (
    projectId: string,
    boardId: string
  ): Promise<CardSummaryView[]> => {
    const response = await apiClient.get<CardSummaryView[]>(
      `/projects/${projectId}/boards/${boardId}/cards`
    );
    return response.data;
  },

  getCardById: async (
    projectId: string,
    boardId: string,
    cardId: string
  ): Promise<CardView> => {
    const response = await apiClient.get<CardView>(
      `/projects/${projectId}/boards/${boardId}/cards/${cardId}`
    );
    return response.data;
  },

  createCard: async (
    projectId: string,
    boardId: string,
    columnId: string,
    payload: CreateCardRequest
  ): Promise<CardView> => {
    const response = await apiClient.post<CardView>(
      `/projects/${projectId}/boards/${boardId}/columns/${columnId}/cards`,
      payload
    );
    return response.data;
  },

  updateCard: async (
    projectId: string,
    boardId: string,
    cardId: string,
    payload: UpdateCardRequest
  ): Promise<CardView> => {
    const response = await apiClient.patch<CardView>(
      `/projects/${projectId}/boards/${boardId}/cards/${cardId}`,
      payload
    );
    return response.data;
  },

  moveCard: async (
    projectId: string,
    boardId: string,
    cardId: string,
    payload: MoveCardRequest
  ): Promise<void> => {
    await apiClient.patch(
      `/projects/${projectId}/boards/${boardId}/cards/${cardId}/move`,
      payload
    );
  },

  deleteCard: async (
    projectId: string,
    boardId: string,
    cardId: string
  ): Promise<void> => {
    await apiClient.delete(
      `/projects/${projectId}/boards/${boardId}/cards/${cardId}`
    );
  },
};
