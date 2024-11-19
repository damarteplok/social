import React from 'react';
import { Box, Typography, Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const UnauthorizedPage: React.FC = () => {
	const navigate = useNavigate();

	const handleSignIn = () => {
		navigate('/sign-in');
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'column',
				alignItems: 'center',
				justifyContent: 'center',
				minHeight: '100vh',
			}}
		>
			<Typography variant='h1' component='div' gutterBottom>
				401
			</Typography>
			<Typography variant='h5' component='div' gutterBottom>
				Page unauthorized
			</Typography>
			<Button variant='text' color='primary' onClick={handleSignIn}>
				Sign In
			</Button>
		</Box>
	);
};

export default UnauthorizedPage;
