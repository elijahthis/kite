import type { CurrencyCode, TransactionDirection } from "../../types";
import { formatCurrency } from "../../utils";

interface AmountTagProps {
	direction: TransactionDirection;
	amount: number;
	currency: CurrencyCode;
}

const AmountTag = ({ direction, amount, currency }: AmountTagProps) => {
	const TEXT_STYLES = {
		DEBIT: "text-gray-900",
		CREDIT: "text-green-600",
	};
	return (
		<span
			className={`font-semibold 
            ${TEXT_STYLES[direction]}`}
		>
			{direction === "CREDIT" ? "+" : "-"} {formatCurrency(amount, currency)}
		</span>
	);
};

export default AmountTag;
