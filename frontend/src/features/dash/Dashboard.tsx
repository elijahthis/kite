import MainLayout from "../../layouts/MainLayout";
import { H2 } from "../../components/headings";
import { formatCurrency } from "../../utils";
import DashBalances from "./DashBalances";
import DashTransactionHistory from "./DashTransactionHistory";

export default function Dashboard() {
	return (
		<MainLayout>
			<div className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8 space-y-8">
				<DashBalances />

				<DashTransactionHistory />
			</div>
		</MainLayout>
	);
}
