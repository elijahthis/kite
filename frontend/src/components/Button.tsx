import type { ReactNode } from "react";

interface ButtonProps {
	type: "submit" | "reset" | "button" | undefined;
	disabled: boolean;
	children: ReactNode;
	onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
}

const Button = ({
	type = "button",
	disabled,
	children,
	onClick,
}: ButtonProps) => {
	return (
		<button
			type={type}
			disabled={disabled}
			className="w-full cursor-pointer bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
			onClick={onClick}
		>
			{children}
		</button>
	);
};

export default Button;
