import type { ReactNode } from "react";

const AuthCard = ({ children }: { children: ReactNode }) => {
	return (
		<div className="max-w-md w-full p-8 bg-white rounded-lg shadow-md border border-gray-100 py-6 px-6">
			{children}
		</div>
	);
};

export default AuthCard;
