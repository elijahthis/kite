export const formatCurrency = (amountInMinor: number, currencyCode: string) => {
	const amountInMajor = amountInMinor / 100;

	return new Intl.NumberFormat("en-US", {
		style: "currency",
		currency: currencyCode,
	}).format(amountInMajor);
};
