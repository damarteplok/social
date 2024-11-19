'use client';
import { SignInPage } from '@toolpad/core/SignInPage';
import { Link, useNavigate } from 'react-router-dom';
import axiosInstance, { setAuthToken } from '../utils/axiosInstance';
import { AxiosError } from 'axios';
import { useDispatch } from 'react-redux';
import { setSession } from '../slices/modules/session/sessionSlice';
import { useNotifications } from '@toolpad/core';

export default function SignIn() {
	const dispatch = useDispatch();
	const navigate = useNavigate();
	const notifications = useNotifications();

	return (
		<SignInPage
			slots={{
				signUpLink: (props) => (
					<Link {...props} to='/sign-up'>
						Sign up
					</Link>
				),
				forgotPasswordLink: (props) => (
					<Link {...props} to='/forgot-password'>
						Forgot password?
					</Link>
				),
			}}
			providers={[{ id: 'credentials', name: 'Email and Password' }]}
			signIn={async (provider, formData, callbackUrl) => {
				try {
					const response = await axiosInstance.post('/authentication/token', {
						email: formData.get('email'),
						password: formData.get('password'),
					});

					notifications.show('Success Sign In', {
						severity: 'success',
						autoHideDuration: 3000,
					});

					const { token, user } = response.data.data;
					setAuthToken(token);
					dispatch(
						setSession({
							user: {
								id: user.id.toString(),
								email: user.email,
								name: user.name,
								image: '',
							},
						})
					);

					navigate(callbackUrl || '/');
				} catch (error) {
					return {
						error:
							error instanceof AxiosError
								? error.response?.data?.error
								: 'An error occurred',
					};
				}
				return {};
			}}
		/>
	);
}
