import { AxiosError } from "axios";
import type { ApiErrorResponse } from "../types";

export const extractApiError = (
	error: unknown,
	fallbackMessage = "An unexpected error occurred",
): string => {
	const axiosError = error as AxiosError<ApiErrorResponse>;

	if (axiosError.isAxiosError && axiosError.response?.data?.message) {
		return axiosError.response.data.message;
	}

	if (error instanceof Error) {
		return error.message;
	}

	return fallbackMessage;
};
