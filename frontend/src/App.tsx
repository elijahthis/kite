import "./App.css";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Login from "./features/auth/login";
import Signup from "./features/auth/signup";
import Dashboard from "./features/dash/Dashboard";
import Deposit from "./features/transactions/Deposit";
// import Signup from './features/auth/Signup';
// import Dashboard from './features/dashboard/Dashboard';

function App() {
	return (
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<Navigate to="/login" replace />} />
				{/* Auth  */}
				<Route path="/login" element={<Login />} />
				<Route path="/signup" element={<Signup />} />
				{/* Protected App Routes */}
				<Route path="/dashboard" element={<Dashboard />} />
				<Route path="/deposit" element={<Deposit />} />
				<Route path="*" element={<Navigate to="/login" replace />} />
			</Routes>
		</BrowserRouter>
	);
}

export default App;
