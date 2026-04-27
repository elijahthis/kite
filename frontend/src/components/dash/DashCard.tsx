interface DashCardProps {
	key: string;
	title: string;
	desc: string;
}

const DashCard = ({ key, title, desc }: DashCardProps) => {
	return (
		<div
			key={key}
			className="p-6 bg-white border border-gray-200 rounded-lg shadow-sm flex flex-col"
		>
			<span className="text-sm font-medium text-gray-500">{title}</span>
			<span className="text-2xl font-bold text-gray-900 mt-2">{desc}</span>
		</div>
	);
};

export default DashCard;
