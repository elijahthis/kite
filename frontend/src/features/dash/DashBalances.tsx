import DashBanner from "../../components/dash/DashBanner";
import DashCard from "../../components/dash/DashCard";
import { H2 } from "../../components/headings";
import { useBalances } from "../../hooks/useDashboardQueries";
import { formatCurrency } from "../../utils";

const DashBalances = () => {
	const {
		data: balances,
		isLoading: loadingBalances,
		isError: errorBalances,
	} = useBalances();
	return (
		<section>
			<H2>Your Wallets</H2>
			{loadingBalances ? (
				<div className="animate-pulse flex space-x-4 mt-4">
					<div className="h-24 bg-gray-200 rounded w-1/4"></div>
					<div className="h-24 bg-gray-200 rounded w-1/4"></div>
				</div>
			) : errorBalances ? (
				<DashBanner>Failed to load balances.</DashBanner>
			) : (
				<div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-4">
					{balances &&
						Object.entries(balances).map(([currency, amount]) => (
							<DashCard
								key={currency}
								title={`${currency} Balance`}
								desc={formatCurrency(amount, currency)}
							/>
						))}
				</div>
			)}
		</section>
	);
};

export default DashBalances;
