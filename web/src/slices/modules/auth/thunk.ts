import { createAsyncThunk } from '@reduxjs/toolkit';
import axiosInstance from '../../../utils/axiosInstance';

export const registerUser = createAsyncThunk(
	'auth/registerUser',
	async (
		{
			username,
			email,
			password,
		}: { username: string; email: string; password: string },
		{ rejectWithValue }
	) => {
		try {
			const response = await axiosInstance.post('/authentication/user', {
				username,
				email,
				password,
			});
			return response.data;
		} catch (error) {
			return rejectWithValue(error);
		}
	}
);
