import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../../api/client";
import MainLayout from "../../layouts/MainLayout";
import { H2 } from "../../components/headings";
import Button from "../../components/Button";
import InputComponent from "../../components/form/InputComponent";
import AuthBanner from "../../components/auth/AuthBanner";
import SelectDropdown from "../../components/form/SelectDropdown";
import type { ApiSuccessResponse, CurrencyCode } from "../../types";
import { extractApiError } from "../../utils/errors";
import { SUPPORTED_CURRENCIES } from "../../utils/constants";

export default function Deposit() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	const [currency, setCurrency] = useState("USD");
	const [amountStr, setAmountStr] = useState("");
	const [successMsg, setSuccessMsg] = useState("");

	interface DepositPayload {
		currency: CurrencyCode;
		amount: number;
		reference: string;
	}

	const depositMutation = useMutation<
		ApiSuccessResponse,
		Error,
		DepositPayload
	>({
		mutationFn: async (payload: {
			currency: string;
			amount: number;
			reference: string;
		}) => {
			const response = await apiClient.post<ApiSuccessResponse>(
				"/deposits",
				payload,
			);
			return response.data;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["balances"] });
			queryClient.invalidateQueries({ queryKey: ["transactions"] });

			setSuccessMsg(`Successfully deposited ${currency} ${amountStr}`);
			setAmountStr("");

			setTimeout(() => navigate("/dashboard"), 2000);
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		setSuccessMsg("");

		const parsedAmount = parseFloat(amountStr);
		if (isNaN(parsedAmount) || parsedAmount <= 0) {
			return;
		}

		const amountInMinorUnits = Math.round(parsedAmount * 100);

		const reference = crypto.randomUUID();

		depositMutation.mutate({
			currency: currency as CurrencyCode,
			amount: amountInMinorUnits,
			reference,
		});
	};

	return (
		<MainLayout>
			<div className="max-w-lg mx-auto py-12 px-4 sm:px-6">
				<div className="bg-white p-8 border border-gray-200 rounded-lg shadow-sm">
					<H2>Fund your Wallet</H2>
					<p className="text-gray-500 text-sm text-center mb-6 mt-1">
						Top up your Wallet today. Receive money instantly into your Kite
						account.
					</p>

					{/* Error State */}
					{depositMutation.isError && (
						<div className="mb-4">
							<AuthBanner>
								{extractApiError(
									depositMutation.error,
									"Failed to process deposit",
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

					<form onSubmit={handleSubmit} className="space-y-5">
						<div>
							<SelectDropdown
								label="Select Currency"
								value={currency}
								setValue={(val: string) => setCurrency(val)}
								valueList={SUPPORTED_CURRENCIES}
							/>
						</div>

						<div>
							<InputComponent
								type="number"
								label="Amount"
								required={true}
								value={amountStr}
								setValue={(val: string) => setAmountStr(val)}
							/>
							<p className="text-xs text-gray-400 mt-1">
								Enter the amount (e.g., 100.50).
							</p>
						</div>

						<div className="pt-2">
							<Button
								type="submit"
								disabled={depositMutation.isPending || !amountStr}
							>
								{depositMutation.isPending
									? "Processing..."
									: "Confirm Deposit"}
							</Button>
						</div>
					</form>
				</div>
			</div>
		</MainLayout>
	);
}
