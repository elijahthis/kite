import type { ReactNode } from "react";

export const H1 = ({ children }: { children: ReactNode }) => {
	return (
		<h1 className="text-3xl font-bold text-center text-gray-900 mb-6">
			{children}
		</h1>
	);
};

export const H2 = ({ children }: { children: ReactNode }) => {
	return (
		<h2 className="text-3xl font-bold text-center text-gray-900 mb-6">
			{children}
		</h2>
	);
};

export const H3 = ({ children }: { children: ReactNode }) => {
	return (
		<h3 className="text-3xl font-bold text-center text-gray-900 mb-6">
			{children}
		</h3>
	);
};

export const H4 = ({ children }: { children: ReactNode }) => {
	return (
		<h4 className="text-3xl font-bold text-center text-gray-900 mb-6">
			{children}
		</h4>
	);
};
