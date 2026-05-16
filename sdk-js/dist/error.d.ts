import type { ErrorCode, ApiErrorResponse } from "./types";
export declare class ApiError extends Error {
    code: ErrorCode;
    details?: Record<string, unknown>;
    constructor(response: ApiErrorResponse);
}
