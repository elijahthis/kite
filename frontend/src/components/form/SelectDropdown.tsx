import type { ChangeEvent } from "react";

interface SelectDropdownProps {
	label: string;
	value: string;
	setValue: (val: string) => void;
	valueList: { label: string; value: string }[];
}

const SelectDropdown = ({
	label,
	value,
	setValue,
	valueList,
}: SelectDropdownProps) => {
	return (
		<label className="block text-sm font-medium text-gray-700 mb-1">
			{label}
			<select
				value={value}
				onChange={(e: ChangeEvent<HTMLSelectElement>) =>
					setValue(e.target.value)
				}
				className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
			>
				{valueList.map((cur) => (
					<option key={cur.value} value={cur.value}>
						{cur.label}
					</option>
				))}
			</select>
		</label>
	);
};

export default SelectDropdown;
