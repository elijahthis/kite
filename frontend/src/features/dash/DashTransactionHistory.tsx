import DashBanner from "../../components/dash/DashBanner";
import { H2 } from "../../components/headings";
import SkeletonLoader from "../../components/SkeletonLoader";
import UITable from "../../components/UITable";
import { useTransactionHistory } from "../../hooks/useDashboardQueries";
import AmountTag from "./AmountTag";
import StatusTag from "./StatusTag";

const DashTransactionHistory = () => {
	const {
		data: transactions,
		isLoading: loadingTxns,
		isError: errorTxns,
	} = useTransactionHistory();

	return (
		<section>
			<H2>Recent Transactions</H2>
			{loadingTxns ? (
				<SkeletonLoader />
			) : errorTxns ? (
				<DashBanner>Failed to load transaction history.</DashBanner>
			) : (
				<UITable
					headings={["TX ID", "Type", "Amount", "Status", "Date"]}
					emptyState={
						<td colSpan={4} className="px-6 py-8 text-center text-gray-500">
							No transactions found.
						</td>
					}
					bodyValues={
						transactions
							? transactions.map((txn) => [
									<span>{txn.transaction_id}</span>,
									<span className="font-bold">{txn.type}</span>,
									<AmountTag
										direction={txn.direction}
										amount={txn.amount}
										currency={txn.currency}
									/>,
									<StatusTag status={txn.status} />,
									<span>{new Date(txn.created_at).toLocaleDateString()}</span>,
								])
							: []
					}
				/>
			)}
		</section>
	);
};

export default DashTransactionHistory;
