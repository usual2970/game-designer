import type { ErrorCode, ApiErrorResponse } from "./types";

export class ApiError extends Error {
  code: ErrorCode;
  details?: Record<string, unknown>;

  constructor(response: ApiErrorResponse) {
    super(response.error);
    this.name = "ApiError";
    this.code = response.code;
    this.details = response.details;
  }
}
