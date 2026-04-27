import DashBanner from "../../components/dash/DashBanner";
import { H2 } from "../../components/headings";
import SkeletonLoader from "../../components/SkeletonLoader";
import { useTransactionHistory } from "../../hooks/useDashboardQueries";
import { formatCurrency } from "../../utils";

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
				<div className="mt-4 bg-white border border-gray-200 rounded-lg shadow-sm overflow-hidden">
					<table className="min-w-full divide-y divide-gray-200">
						<thead className="bg-gray-50">
							<tr>
								<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
									Date
								</th>
								<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
									Type
								</th>
								<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
									Amount
								</th>
								<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
									Status
								</th>
							</tr>
						</thead>
						<tbody className="bg-white divide-y divide-gray-200">
							{transactions?.length === 0 ? (
								<tr>
									<td
										colSpan={4}
										className="px-6 py-8 text-center text-gray-500"
									>
										No transactions found.
									</td>
								</tr>
							) : (
								transactions?.map((txn) => (
									<tr key={txn.transaction_id}>
										<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{new Date(txn.created_at).toLocaleDateString()}
										</td>
										<td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
											{txn.type}
										</td>
										<td
											className={`px-6 py-4 whitespace-nowrap text-sm font-semibold ${txn.direction === "CREDIT" ? "text-green-600" : "text-gray-900"}`}
										>
											{txn.direction === "CREDIT" ? "+" : "-"}{" "}
											{formatCurrency(txn.amount, txn.currency)}
										</td>
										<td className="px-6 py-4 whitespace-nowrap text-sm">
											<span
												className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                                                        ${
																													txn.status ===
																													"SUCCESS"
																														? "bg-green-100 text-green-800"
																														: txn.status ===
																																	"PENDING" ||
																															  txn.status ===
																																	"PROCESSING"
																															? "bg-yellow-100 text-yellow-800"
																															: "bg-red-100 text-red-800"
																												}`}
											>
												{txn.status}
											</span>
										</td>
									</tr>
								))
							)}
						</tbody>
					</table>
				</div>
			)}
		</section>
	);
};

export default DashTransactionHistory;
