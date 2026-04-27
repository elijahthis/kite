import type { TransactionStatus } from "../../types";

interface StatusTagProps {
	status: TransactionStatus;
}

const StatusTag = ({ status }: StatusTagProps) => {
	const TEXT_STYLES = {
		SUCCESS: "bg-green-100 text-green-800",
		PENDING: "bg-yellow-100 text-yellow-800",
		PROCESSING: "bg-yellow-100 text-yellow-800",
		FAILED: "bg-red-100 text-red-800",
	};
	return (
		<span
			className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
            ${TEXT_STYLES[status]}`}
		>
			{status}
		</span>
	);
};

export default StatusTag;
