import "./App.css";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Login from "./features/auth/login";
import Signup from "./features/auth/signup";
import Dashboard from "./features/dash/Dashboard";
import Deposit from "./features/transactions/Deposit";
import Conversion from "./features/transactions/Conversions";
import Payout from "./features/transactions/Payout";

function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Navigate to="/login" replace />} />
				{/* Auth  */}
				<Route path="/login" element={<Login />} />
				<Route path="/signup" element={<Signup />} />

				{/* Protected Routes */}
				<Route path="/dashboard" element={<Dashboard />} />
				<Route path="/deposit" element={<Deposit />} />
				<Route path="/conversion" element={<Conversion />} />
				<Route path="/payout" element={<Payout />} />

				<Route path="*" element={<Navigate to="/login" replace />} />
			</Routes>
		</BrowserRouter>
	);
}

export default App;
