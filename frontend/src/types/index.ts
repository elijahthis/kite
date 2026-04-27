// --- Core Domain Enums ---
export type CurrencyCode = "USD" | "GBP" | "EUR" | "NGN" | "KES";
export type TransactionType = "DEPOSIT" | "PAYOUT" | "CONVERSION" | "REVERSAL";
export type TransactionStatus = "PENDING" | "PROCESSING" | "SUCCESS" | "FAILED";
export type TransactionDirection = "CREDIT" | "DEBIT";

// --- Standard API Error ---
export interface ApiErrorResponse {
	error: string;
	message: string;
}

// --- Standard Success Response ---
export interface ApiSuccessResponse {
	status: string;
	message?: string;
}

// --- Models ---
export interface Transaction {
	transaction_id: string;
	type: TransactionType;
	status: TransactionStatus;
	amount: number;
	direction: TransactionDirection;
	currency: CurrencyCode;
	reference: string;
	created_at: string;
}

export interface Quote {
	quote_id: string;
	exchange_rate: number;
	amount_in: number;
	amount_out: number;
	expires_at: string;
}

// --- API Responses ---
export interface BalancesResponse extends ApiSuccessResponse {
	balances: Record<CurrencyCode, number>;
}

export interface HistoryResponse extends ApiSuccessResponse {
	transactions: Transaction[];
	page: number;
	limit: number;
}
