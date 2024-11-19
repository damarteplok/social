import { createAsyncThunk } from '@reduxjs/toolkit';
import { resetSession } from './sessionSlice';
import axiosInstance from '../../../utils/axiosInstance';

export const checkAuth = createAsyncThunk(
	'session/checkAuth',
	async (_, { dispatch, rejectWithValue }) => {
		const token = localStorage.getItem('token');
		if (token) {
			try {
				const response = await axiosInstance.get('/authentication/user');
				return response.data.data;
			} catch (error) {
				dispatch(resetSession());
				localStorage.removeItem('token');
				rejectWithValue(error);
			}
		} else {
			dispatch(resetSession());
			localStorage.removeItem('token');
			rejectWithValue('No token found');
		}
	}
);
