export class ApiError extends Error {
    constructor(response) {
        super(response.error);
        this.name = "ApiError";
        this.code = response.code;
        this.details = response.details;
    }
}
