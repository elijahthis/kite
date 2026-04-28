import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { apiClient } from "../../api/client";
import { H2 } from "../../components/headings";
import AuthBanner from "../../components/auth/AuthBanner";
import InputComponent from "../../components/form/InputComponent";
import AuthCard from "../../components/auth/AuthCard";
import AuthLayout from "../../layouts/AuthLayout";
import Button from "../../components/Button";
import { extractApiError } from "../../utils/errors";

export default function Signup() {
	const navigate = useNavigate();
	const [firstName, setFirstName] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const signupMutation = useMutation({
		mutationFn: async () => {
			const response = await apiClient.post("/auth/register", {
				first_name: firstName,
				email,
				password,
			});
			return response.data;
		},
		onSuccess: () => {
			navigate("/login", {
				state: { message: "Account created successfully. Please log in." },
			});
		},
	});

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		signupMutation.mutate();
	};

	return (
		<AuthLayout>
			<AuthCard>
				<H2>Create a Kite Account</H2>

				{/* Error State */}
				{signupMutation.isError && (
					<AuthBanner>
						{extractApiError(signupMutation.error, "Failed to create account")}
					</AuthBanner>
				)}

				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<InputComponent
							type="text"
							label="First Name"
							required={true}
							value={firstName}
							setValue={(val: string) => {
								setFirstName(val);
							}}
						/>
					</div>
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

					<Button type="submit" disabled={signupMutation.isPending}>
						{signupMutation.isPending ? "Creating Account..." : "Sign Up"}
					</Button>
				</form>

				<p className="mt-4 text-center text-sm text-gray-600">
					Already have an account?{" "}
					<Link to="/login" className="text-blue-600 hover:underline">
						Log in
					</Link>
				</p>
			</AuthCard>
		</AuthLayout>
	);
}
