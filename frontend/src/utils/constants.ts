import type { CurrencyCode } from "../types";

export const SUPPORTED_CURRENCIES: { label: string; value: CurrencyCode }[] = [
	{ label: "USD", value: "USD" },
	{ label: "GBP", value: "GBP" },
	{ label: "EUR", value: "EUR" },
	{ label: "NGN", value: "NGN" },
	{ label: "KES", value: "KES" },
];
