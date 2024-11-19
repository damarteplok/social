import React from 'react';
import {
	Box,
	Button,
	Container,
	TextField,
	Typography,
	Alert,
	Card,
	CardContent,
	Link,
	Paper,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';
import { useFormik } from 'formik';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../slices/store/rootReducer';
import { registerUser } from '../slices/modules/auth/thunk';

interface SignUpProps {
	onSuccess: () => void;
}

const SignUp: React.FC<SignUpProps> = ({ onSuccess }) => {
	const navigate = useNavigate();
	const dispatch = useDispatch<any>();
	const registerUserState = useSelector(
		(state: RootState) => state.registerUser
	);

	const formik = useFormik({
		enableReinitialize: true,
		initialValues: {
			username: '',
			email: '',
			password: '',
		},
		validationSchema: Yup.object({
			username: Yup.string().required('Required'),
			email: Yup.string().email('Invalid email address').required('Required'),
			password: Yup.string().required('Required'),
		}),
		onSubmit: async (values, { setSubmitting }) => {
			try {
				setSubmitting(true);
				dispatch(
					registerUser({
						username: values.username,
						email: values.email,
						password: values.password,
					})
				);
				setSubmitting(false);
				onSuccess();
			} catch (err) {
				setSubmitting(false);
			}
		},
	});

	const handleSignInClick = () => {
		navigate('/sign-in');
	};

	return (
		<Container maxWidth='sm'>
			<Box
				sx={{
					mt: 8,
					display: 'flex',
					flexDirection: 'column',
					alignItems: 'center',
					justifyContent: 'center',
					minHeight: '100vh',
				}}
			>
				<Paper elevation={8}>
					<Card sx={{ width: '100%', maxWidth: 400 }}>
						<CardContent>
							<Typography component='h1' variant='h5' align='center'>
								Sign Up
							</Typography>
							<Box
								component='form'
								onSubmit={formik.handleSubmit}
								sx={{ mt: 1 }}
							>
								{registerUserState.errorMessage && (
									<Alert severity='error'>
										{registerUserState.errorMessage}
									</Alert>
								)}
								{registerUserState.success && (
									<Alert severity='success'>User registered successfully</Alert>
								)}
								<TextField
									margin='normal'
									required
									fullWidth
									id='username'
									label='User Name'
									name='username'
									autoComplete='username'
									autoFocus
									value={formik.values.username}
									onChange={formik.handleChange}
									onBlur={formik.handleBlur}
									error={
										formik.touched.username && Boolean(formik.errors.username)
									}
									helperText={formik.touched.username && formik.errors.username}
								/>
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
									error={formik.touched.email && Boolean(formik.errors.email)}
									helperText={formik.touched.email && formik.errors.email}
								/>
								<TextField
									margin='normal'
									required
									fullWidth
									name='password'
									label='Password'
									type='password'
									id='password'
									autoComplete='current-password'
									value={formik.values.password}
									onChange={formik.handleChange}
									onBlur={formik.handleBlur}
									error={
										formik.touched.password && Boolean(formik.errors.password)
									}
									helperText={formik.touched.password && formik.errors.password}
								/>
								<Button
									type='submit'
									fullWidth
									variant='contained'
									sx={{ mt: 3, mb: 2 }}
									disabled={formik.isSubmitting}
								>
									Sign Up
								</Button>
								<Box textAlign='center'>
									<Link
										component='button'
										variant='body2'
										onClick={handleSignInClick}
									>
										Already have an account? Sign in
									</Link>
								</Box>
							</Box>
						</CardContent>
					</Card>
				</Paper>
			</Box>
		</Container>
	);
};

export default SignUp;
