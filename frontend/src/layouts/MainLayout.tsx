import type { ReactNode } from "react";

const MainLayout = ({ children }: { children: ReactNode }) => {
	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-50 ">
			{children}
		</div>
	);
};

export default MainLayout;
