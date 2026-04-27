import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../../api/client";
import MainLayout from "../../layouts/MainLayout";
import { H2 } from "../../components/headings";
import Button from "../../components/Button";
import InputComponent from "../../components/form/InputComponent";
import SelectDropdown from "../../components/form/SelectDropdown";
import AuthBanner from "../../components/auth/AuthBanner";
import { extractApiError } from "../../utils/errors";
import type { CurrencyCode, ApiSuccessResponse } from "../../types";
import { SUPPORTED_CURRENCIES } from "../../utils/constants";

interface PayoutPayload {
	source_currency: CurrencyCode;
	amount: number;
	account_number: string;
	bank_code: string;
	account_name: string;
}

interface PayoutResponse extends ApiSuccessResponse {
	transaction_id: string;
}

export default function Payout() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	// Form State
	const [currency, setCurrency] = useState<CurrencyCode>("USD");
	const [amountStr, setAmountStr] = useState("");
	const [accountNumber, setAccountNumber] = useState("");
	const [bankCode, setBankCode] = useState("");
	const [accountName, setAccountName] = useState("");

	const [successMsg, setSuccessMsg] = useState("");

	const payoutMutation = useMutation<PayoutResponse, Error, PayoutPayload>({
		mutationFn: async (payload) => {
			const response = await apiClient.post<PayoutResponse>(
				"/payouts",
				payload,
			);
			return response.data;
		},
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: ["balances"] });
			queryClient.invalidateQueries({ queryKey: ["transactions"] });

			setSuccessMsg(
				`Payout initiated successfully! Transaction ID: ${data.transaction_id}`,
			);

			// Reset form
			setAmountStr("");
			setAccountNumber("");
			setBankCode("");
			setAccountName("");

			setTimeout(() => navigate("/dashboard"), 3000);
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		setSuccessMsg("");

		const parsedAmount = parseFloat(amountStr);
		if (isNaN(parsedAmount) || parsedAmount <= 0) return;

		const amountInMinorUnits = Math.round(parsedAmount * 100);

		payoutMutation.mutate({
			source_currency: currency,
			amount: amountInMinorUnits,
			account_number: accountNumber,
			bank_code: bankCode,
			account_name: accountName,
		});
	};

	return (
		<MainLayout>
			<div className="max-w-lg mx-auto py-12 px-4 sm:px-6">
				<div className="bg-white p-8 border border-gray-200 rounded-lg shadow-sm">
					<H2>Withdraw Funds</H2>
					<p className="text-gray-500 text-sm mb-6 mt-1">
						Send money from your Kite wallet to a bank account.
					</p>

					{/* Error State */}
					{payoutMutation.isError && (
						<div className="mb-4">
							<AuthBanner>
								{extractApiError(
									payoutMutation.error,
									"Failed to initiate payout",
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
						<div className="grid grid-cols-2 gap-4">
							<SelectDropdown
								label="Currency"
								value={currency}
								setValue={(val) => setCurrency(val as CurrencyCode)}
								valueList={SUPPORTED_CURRENCIES}
							/>
							<InputComponent
								type="number"
								label="Amount"
								required={true}
								value={amountStr}
								setValue={(val: string) => setAmountStr(val)}
							/>
						</div>

						<InputComponent
							type="text"
							label="Recipient Account Name"
							required={true}
							value={accountName}
							setValue={(val: string) => setAccountName(val)}
						/>

						<div className="grid grid-cols-2 gap-4">
							<InputComponent
								type="text"
								label="Account Number"
								required={true}
								value={accountNumber}
								setValue={(val: string) => setAccountNumber(val)}
							/>
							<InputComponent
								type="text"
								label="Bank Code / Routing"
								required={true}
								value={bankCode}
								setValue={(val: string) => setBankCode(val)}
							/>
						</div>

						<div className="pt-2">
							<Button
								type="submit"
								disabled={
									payoutMutation.isPending ||
									!amountStr ||
									!accountNumber ||
									!bankCode ||
									!accountName
								}
							>
								{payoutMutation.isPending ? "Initiating..." : "Send Payout"}
							</Button>
						</div>
					</form>
				</div>
			</div>
		</MainLayout>
	);
}
