import { useQuery } from "@tanstack/react-query";
import { apiClient } from "../api/client";

export interface Transaction {
	transaction_id: string;
	type: string;
	status: string;
	amount: number;
	direction: "CREDIT" | "DEBIT";
	currency: string;
	reference: string;
	created_at: string;
}

export const useBalances = () => {
	return useQuery({
		queryKey: ["balances"],
		queryFn: async () => {
			const response = await apiClient.get("/balances");
			return response.data.balances as Record<string, number>;
		},
	});
};

export const useTransactionHistory = (page = 1, limit = 20) => {
	return useQuery({
		queryKey: ["transactions", page, limit],
		queryFn: async () => {
			const response = await apiClient.get(
				`/transactions?page=${page}&limit=${limit}`,
			);
			return response.data.transactions as Transaction[];
		},
	});
};
