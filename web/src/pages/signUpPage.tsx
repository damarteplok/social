import React from 'react';
import SignUp from '../components/page/SignUp';

const SignUpPage: React.FC = () => {
	const handleSuccess = () => {};

	return <SignUp onSuccess={handleSuccess} />;
};

export default SignUpPage;
