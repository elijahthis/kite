import { Link } from "react-router-dom";
import { H2 } from "../../components/headings";

const DashQuickActions = () => {
	type ActionItem = { label: string; route: string };
	const ACTIONS: ActionItem[] = [
		{
			label: "+ Deposit",
			route: "/deposit",
		},
		{
			label: "Convert FX",
			route: "/conversion",
		},
		{
			label: "Payout",
			route: "/payout",
		},
	];

	return (
		<div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
			<H2>Dashboard</H2>
			<div className="flex space-x-3">
				{ACTIONS.map((item: ActionItem, ind: number) => (
					<Link
						to={item.route}
						className={`${ind == 0 ? "bg-blue-600 text-white hover:bg-blue-700" : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"}
                         px-4 py-2 rounded-md  transition-colors text-sm font-medium shadow-sm`}
					>
						{item.label}
					</Link>
				))}
			</div>
		</div>
	);
};

export default DashQuickActions;
