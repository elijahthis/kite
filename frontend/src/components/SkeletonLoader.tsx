const SkeletonLoader = ({
	skeletonClass = "h-10 bg-gray-200 rounded w-full",
}: {
	skeletonClass?: string;
}) => {
	return (
		<div className="mt-4 space-y-3 animate-pulse">
			<div className={skeletonClass}></div>
			<div className={skeletonClass}></div>
		</div>
	);
};

export default SkeletonLoader;
