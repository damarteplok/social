import React, { useState } from 'react';
import {
	Box,
	Button,
	Container,
	TextField,
	Typography,
	Alert,
} from '@mui/material';
import axiosInstance from '../utils/axiosInstance';

interface ForgotPasswordProps {
	onSuccess: () => void;
	onBack: () => void;
}

const ForgotPassword: React.FC<ForgotPasswordProps> = ({
	onSuccess,
	onBack,
}) => {
	const [email, setEmail] = useState('');
	const [error, setError] = useState<string | null>(null);
	const [success, setSuccess] = useState(false);

	const handleSubmit = async (event: React.FormEvent) => {
		event.preventDefault();
		setError(null);
		setSuccess(false);

		try {
			await axiosInstance.post('/authentication/forgot-password', { email });
			setSuccess(true);
			onSuccess();
		} catch (err) {
			setError('Failed to send reset password email. Please try again.');
		}
	};

	return (
		<Container maxWidth='sm'>
			<Box
				sx={{
					mt: 8,
					display: 'flex',
					flexDirection: 'column',
					alignItems: 'center',
				}}
			>
				<Typography component='h1' variant='h5'>
					Forgot Password
				</Typography>
				<Box component='form' onSubmit={handleSubmit} sx={{ mt: 1 }}>
					<TextField
						margin='normal'
						required
						fullWidth
						id='email'
						label='Email Address'
						name='email'
						autoComplete='email'
						autoFocus
						value={email}
						onChange={(e) => setEmail(e.target.value)}
						size="small"
					/>
					{error && <Alert severity='error'>{error}</Alert>}
					{success && (
						<Alert severity='success'>
							Reset password email sent successfully.
						</Alert>
					)}
					<Button
						type='submit'
						fullWidth
						variant='contained'
						sx={{ mt: 3, mb: 2 }}
					>
						Send Reset Link
					</Button>
					<Button
						type='button'
						fullWidth
						variant='text'
						sx={{ mt: 1, mb: 2 }}
						onClick={() => onBack()}
					>
						Back to Sign In
					</Button>
				</Box>
			</Box>
		</Container>
	);
};

export default ForgotPassword;
