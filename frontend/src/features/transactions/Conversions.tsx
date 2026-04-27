import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../../api/client";
import MainLayout from "../../layouts/MainLayout";
import { H2 } from "../../components/headings";
import Button from "../../components/Button";
import InputComponent from "../../components/form/InputComponent";
import SelectDropdown from "../../components/form/SelectDropdown"; // Your new component
import AuthBanner from "../../components/auth/AuthBanner";
import { extractApiError } from "../../utils/errors";
import type { Quote, CurrencyCode, ApiSuccessResponse } from "../../types";
import ConversionsQuoteReview from "./ConversionsQuoteReview";
import { SUPPORTED_CURRENCIES } from "../../utils/constants";

export default function Conversion() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	// Form State
	const [sourceCurrency, setSourceCurrency] = useState<CurrencyCode>("USD");
	const [targetCurrency, setTargetCurrency] = useState<CurrencyCode>("NGN");
	const [amountStr, setAmountStr] = useState("");

	// Quote State
	const [quote, setQuote] = useState<Quote | null>(null);
	const [successMsg, setSuccessMsg] = useState("");

	// Step 1: Fetch Quote Mutation
	interface QuotePayload {
		source_currency: CurrencyCode;
		target_currency: CurrencyCode;
		amount_in: number;
	}
	const quoteMutation = useMutation<Quote, Error, QuotePayload>({
		mutationFn: async (payload: {
			source_currency: string;
			target_currency: string;
			amount_in: number;
		}) => {
			const response = await apiClient.post<Quote>(
				"/conversions/quote",
				payload,
			);
			return response.data;
		},
		onSuccess: (data) => {
			setQuote(data);
		},
	});

	// Step 2: Execute Quote Mutation
	const executeMutation = useMutation<ApiSuccessResponse, Error, string>({
		mutationFn: async (quote_id: string) => {
			const response = await apiClient.post<ApiSuccessResponse>(
				"/conversions/execute",
				{
					quote_id,
				},
			);
			return response.data;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["balances"] });
			queryClient.invalidateQueries({ queryKey: ["transactions"] });

			setSuccessMsg("Conversion executed successfully!");
			setQuote(null);
			setAmountStr("");

			setTimeout(() => navigate("/dashboard"), 2000);
		},
	});

	const handleGetQuote = (e: React.FormEvent) => {
		e.preventDefault();
		setSuccessMsg("");
		setQuote(null);

		if (sourceCurrency === targetCurrency) return;

		const parsedAmount = parseFloat(amountStr);
		if (isNaN(parsedAmount) || parsedAmount <= 0) return;

		const amountInMinorUnits = Math.round(parsedAmount * 100);

		quoteMutation.mutate({
			source_currency: sourceCurrency,
			target_currency: targetCurrency,
			amount_in: amountInMinorUnits,
		});
	};

	const handleExecute = (event: React.MouseEvent<HTMLButtonElement>) => {
		if (!quote?.quote_id) return;
		executeMutation.mutate(quote.quote_id);
	};

	return (
		<MainLayout>
			<div className="max-w-lg mx-auto py-12 px-4 sm:px-6">
				<div className="bg-white p-8 border border-gray-200 rounded-lg shadow-sm">
					<H2>Convert Currencies</H2>
					<p className="text-gray-500 text-sm mb-6 mt-1">
						Exchange funds instantly across your wallets.
					</p>

					{/* Error States */}
					{(quoteMutation.isError || executeMutation.isError) && (
						<div className="mb-4">
							<AuthBanner>
								{quoteMutation.isError
									? extractApiError(quoteMutation.error, "Failed to get quote")
									: extractApiError(
											executeMutation.error,
											"Failed to execute conversion",
										)}
							</AuthBanner>
						</div>
					)}

					{/* Success State */}
					{successMsg && (
						<div className="mb-4 p-3 bg-green-50 text-green-700 rounded-md text-sm font-medium">
							{successMsg}
						</div>
					)}

					<form onSubmit={handleGetQuote} className="space-y-5">
						<div className="grid grid-cols-2 gap-4">
							<SelectDropdown
								label="From"
								value={sourceCurrency}
								setValue={setSourceCurrency}
								valueList={SUPPORTED_CURRENCIES}
							/>
							<SelectDropdown
								label="To"
								value={targetCurrency}
								setValue={setTargetCurrency}
								valueList={SUPPORTED_CURRENCIES}
							/>
						</div>

						{sourceCurrency === targetCurrency && (
							<p className="text-xs text-red-500">
								Source and target currencies must be different.
							</p>
						)}

						<InputComponent
							type="number"
							label="Amount to Convert"
							required={true}
							value={amountStr}
							setValue={(val: string) => {
								setAmountStr(val);
								setQuote(null);
							}}
						/>

						{!quote && (
							<div className="pt-2">
								<Button
									type="submit"
									disabled={
										quoteMutation.isPending ||
										!amountStr ||
										sourceCurrency === targetCurrency
									}
								>
									{quoteMutation.isPending ? "Fetching Rate..." : "Get Quote"}
								</Button>
							</div>
						)}
					</form>

					{quote && !successMsg && (
						<ConversionsQuoteReview
							sourceCurrency={sourceCurrency}
							targetCurrency={targetCurrency}
							quote={quote}
						>
							<Button
								type="button"
								onClick={handleExecute}
								disabled={executeMutation.isPending}
							>
								{executeMutation.isPending
									? "Executing..."
									: "Confirm & Convert"}
							</Button>
						</ConversionsQuoteReview>
					)}
				</div>
			</div>
		</MainLayout>
	);
}
