import React from 'react';
import {
	Box,
	Button,
	Container,
	TextField,
	Typography,
	Alert,
} from '@mui/material';
import * as Yup from 'yup';
import { useFormik } from 'formik';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../slices/store/rootReducer';
import { forgotPassword } from '../slices/modules/forgotPassword/thunk';
interface ForgotPasswordProps {
	onSuccess: () => void;
	onBack: () => void;
}

const ForgotPassword: React.FC<ForgotPasswordProps> = ({
	onSuccess,
	onBack,
}) => {
	const dispatch = useDispatch<any>();
	const { success, error } = useSelector(
		(state: RootState) => state.forgotPassword
	);

	const formik = useFormik({
		enableReinitialize: true,
		initialValues: {
			email: '',
		},
		validationSchema: Yup.object({
			email: Yup.string().email('Invalid email address').required('Required'),
		}),
		onSubmit: async (values, { setSubmitting }) => {
			try {
				setSubmitting(true);
				dispatch(forgotPassword(values.email));
				setSubmitting(false);
				onSuccess();
			} catch (err) {
				setSubmitting(false);
			}
		},
	});

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
				<Box
					component='form'
					onSubmit={formik.handleSubmit}
					sx={{ mt: 1, minWidth: 350 }}
				>
					{error && (
						<Alert severity='error' sx={{ width: '100%' }}>
							{error}
						</Alert>
					)}
					{success && (
						<Alert sx={{ width: '100%' }} severity='success'>
							Reset password email sent successfully.
						</Alert>
					)}
					<TextField
						margin='normal'
						required
						fullWidth
						id='email'
						label='Email Address'
						name='email'
						autoComplete='email'
						value={formik.values.email}
						onChange={formik.handleChange}
						onBlur={formik.handleBlur}
						size='small'
						error={formik.touched.email && Boolean(formik.errors.email)}
						helperText={formik.touched.email && formik.errors.email}
					/>
					<Button
						type='submit'
						fullWidth
						variant='contained'
						sx={{ mt: 3, mb: 2 }}
						disabled={formik.isSubmitting}
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
