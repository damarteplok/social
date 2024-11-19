import React from 'react';
import { Box, Typography, Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const NotFoundPage: React.FC = () => {
	const navigate = useNavigate();

	const handleGoHome = () => {
		navigate('/');
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
				404
			</Typography>
			<Typography variant='h5' component='div' gutterBottom>
				Page Not Found
			</Typography>
			<Button variant='text' color='primary' onClick={handleGoHome}>
				Take me home
			</Button>
		</Box>
	);
};

export default NotFoundPage;
