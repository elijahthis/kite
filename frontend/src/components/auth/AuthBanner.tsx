import type { ReactNode } from "react";

const AuthBanner = ({ children }: { children: ReactNode }) => {
	return (
		<div className="mb-4 p-3 bg-red-50 text-red-700 rounded-md text-sm">
			{children}
		</div>
	);
};

export default AuthBanner;
