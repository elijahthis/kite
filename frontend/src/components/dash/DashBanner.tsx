import type { ReactNode } from "react";

const DashBanner = ({ children }: { children: ReactNode }) => {
	return (
		<div className="mt-4 p-4 text-red-700 bg-red-50 rounded-md">{children}</div>
	);
};

export default DashBanner;
