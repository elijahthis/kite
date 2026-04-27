import type { ChangeEvent } from "react";

interface InputComponent {
	value: string;
	setValue: (val: string) => void;
	label: string;
	type: string;
	required: boolean;
}

const InputComponent = ({
	value,
	setValue,
	label,
	type,
	required,
}: InputComponent) => {
	return (
		<label className="block text-sm font-medium text-gray-700 mb-1">
			{label}
			<input
				className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				value={value}
				onChange={(e: ChangeEvent<HTMLInputElement>) => {
					setValue(e.target.value);
				}}
				type={type}
				required={required}
			/>
		</label>
	);
};

export default InputComponent;
