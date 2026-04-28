import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { apiClient } from "../../api/client";
import InputComponent from "../../components/form/InputComponent";
import AuthCard from "../../components/auth/AuthCard";
import AuthLayout from "../../layouts/AuthLayout";
import { H2 } from "../../components/headings";
import AuthBanner from "../../components/auth/AuthBanner";
import Button from "../../components/Button";
import { extractApiError } from "../../utils/errors";

export default function Login() {
	const navigate = useNavigate();
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const loginMutation = useMutation({
		mutationFn: async () => {
			const response = await apiClient.post("/auth/login", { email, password });
			return response.data;
		},
		onSuccess: () => {
			navigate("/dashboard");
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		loginMutation.mutate();
	};

	return (
		<AuthLayout>
			<AuthCard>
				<H2>Log in to Kite by Grey</H2>

				{/* Error State */}
				{loginMutation.isError && (
					<AuthBanner>
						{extractApiError(loginMutation.error, "Failed to login")}
					</AuthBanner>
				)}

				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<InputComponent
							type="email"
							label="Email"
							required={true}
							value={email}
							setValue={(val: string) => {
								setEmail(val);
							}}
						/>
					</div>
					<div>
						<InputComponent
							type="password"
							label="Password"
							required={true}
							value={password}
							setValue={(val: string) => {
								setPassword(val);
							}}
						/>
					</div>

					<Button type="submit" disabled={loginMutation.isPending}>
						{loginMutation.isPending ? "Logging in..." : "Log In"}
					</Button>
				</form>

				<p className="mt-4 text-center text-sm text-gray-600">
					Don't have an account?{" "}
					<Link to="/signup" className="text-blue-600 hover:underline">
						Sign up
					</Link>
				</p>
			</AuthCard>
		</AuthLayout>
	);
}
