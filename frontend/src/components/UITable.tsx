import type { ReactNode } from "react";

interface UITableProps {
	headings: string[];
	emptyState: ReactNode;
	bodyValues: ReactNode[][];
}

const UITable = ({ headings, emptyState, bodyValues }: UITableProps) => {
	return (
		<table className="min-w-full divide-y divide-gray-200">
			<thead className="bg-gray-50">
				<tr>
					{headings.map((item) => (
						<th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
							{item}
						</th>
					))}
				</tr>
			</thead>
			<tbody className="bg-white divide-y divide-gray-200">
				{bodyValues?.length === 0 ? (
					<tr>{emptyState}</tr>
				) : (
					bodyValues.map((row) => (
						<tr>
							{row.map((cell) => (
								<td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
									{cell}
								</td>
							))}
						</tr>
					))
				)}
			</tbody>
		</table>
	);
};

export default UITable;
