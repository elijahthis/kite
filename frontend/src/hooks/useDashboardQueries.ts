import { useQuery } from "@tanstack/react-query";
import { apiClient } from "../api/client";
import type { BalancesResponse, HistoryResponse, Transaction } from "../types";

export const useBalances = () => {
	return useQuery({
		queryKey: ["balances"],
		queryFn: async () => {
			const response = await apiClient.get<BalancesResponse>("/balances");
			return response.data.balances as Record<string, number>;
		},
	});
};

export const useTransactionHistory = (page = 1, limit = 20) => {
	return useQuery({
		queryKey: ["transactions", page, limit],
		queryFn: async () => {
			const response = await apiClient.get<HistoryResponse>(
				`/transactions?page=${page}&limit=${limit}`,
			);
			return response.data.transactions as Transaction[];
		},
	});
};
