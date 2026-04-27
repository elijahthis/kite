import type { ReactNode } from "react";
import type { Quote } from "../../types";
import { formatCurrency } from "../../utils";

interface ConversionsQuoteReviewProps {
	sourceCurrency: string;
	targetCurrency: string;
	quote: Quote;
	children: ReactNode;
}

const ConversionsQuoteReview = ({
	sourceCurrency,
	targetCurrency,
	quote,
	children,
}: ConversionsQuoteReviewProps) => {
	return (
		<div className="mt-6 pt-6 border-t border-gray-200">
			<h3 className="text-sm font-semibold text-gray-900 mb-4">
				Review Conversion
			</h3>

			<div className="bg-gray-50 p-4 rounded-md space-y-3 mb-6">
				<div className="flex justify-between text-sm">
					<span className="text-gray-500">Exchange Rate</span>
					<span className="font-medium text-gray-900">
						1 {sourceCurrency} = {quote.exchange_rate.toFixed(4)}{" "}
						{targetCurrency}
					</span>
				</div>
				<div className="flex justify-between text-sm">
					<span className="text-gray-500">You Spend</span>
					<span className="font-medium text-gray-900">
						{formatCurrency(quote.amount_in, sourceCurrency)}
					</span>
				</div>
				<div className="flex justify-between text-lg font-bold border-t border-gray-200 pt-3 mt-3">
					<span className="text-gray-900">You Receive</span>
					<span className="text-green-600">
						{formatCurrency(quote.amount_out, targetCurrency)}
					</span>
				</div>
			</div>
			{children}

			<p className="text-center text-xs text-gray-400 mt-3">
				Quote expires in 60 seconds.
			</p>
		</div>
	);
};

export default ConversionsQuoteReview;
