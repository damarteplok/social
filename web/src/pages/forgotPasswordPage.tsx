import React from 'react';
import ForgotPassword from '../components/page/ForgotPassword';
import { useNavigate } from 'react-router-dom';

const ForgotPasswordPage: React.FC = () => {
	const navigate = useNavigate();
	const handleSuccess = () => {};

	const handleBack = () => {
		navigate('/sign-in');
	};

	return <ForgotPassword onSuccess={handleSuccess} onBack={handleBack} />;
};

export default ForgotPasswordPage;
