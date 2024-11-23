import { createAsyncThunk } from '@reduxjs/toolkit';
import axiosInstance from '../../../utils/axiosInstance';

export const forgotPassword = createAsyncThunk(
	'forgotPassword/forgotPassword',
	async (email: string, { rejectWithValue }) => {
		try {
			const response = await axiosInstance.post('/auth/forgot-password', {
				email,
			});

			const data = response.data;

			if (response.status === 200) {
				return data;
			}

			return rejectWithValue(data.message);
		} catch (error) {
			return rejectWithValue(error);
		}
	}
);
