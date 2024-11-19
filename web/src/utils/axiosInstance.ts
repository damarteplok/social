import axios from 'axios';
import NProgress from 'nprogress';
import { store } from '../slices/store/store';
import { setSession } from '../slices/modules/session/sessionSlice';

const baseURL = import.meta.env.VITE_API_URL as string;

const axiosInstance = axios.create({
	baseURL,
});

axiosInstance.interceptors.request.use((config) => {
	NProgress.configure({ showSpinner: false });
	NProgress.start();
	const token = localStorage.getItem('token');
	if (token) {
		config.headers.Authorization = `Bearer ${token}`;
	}
	return config;
});

axiosInstance.interceptors.response.use(
	function (response) {
		NProgress.done();
		return response;
	},
	function (error) {
		let message;
		if (axios.isAxiosError(error) && error.response) {
			const status = error.response.status;
			if (status === 401) {
				message = 'Invalid credentials';
				store.dispatch(setSession(null));
				setAuthToken(null);
				window.location.href = '/unauthorized';
			} else {
				message =
					error.response.data?.error ||
					error.response.data?.message ||
					'An error occurred';
			}
		} else if (error.request) {
			message = 'No response received from the server';
		} else {
			message = error.message || 'Request configuration error';
		}
		NProgress.done();
		return Promise.reject(message);
	}
);

export const setAuthToken = (token: string | null) => {
	if (token) {
		localStorage.setItem('token', token);
		axiosInstance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
	} else {
		localStorage.removeItem('token');
		delete axiosInstance.defaults.headers.common['Authorization'];
	}
};

export default axiosInstance;
