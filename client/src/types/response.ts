export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface PaginationResult<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface ApiErrorResponse {
  code: any;
  message: any;
  error?: any;
}
