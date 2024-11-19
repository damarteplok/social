import { createAsyncThunk } from '@reduxjs/toolkit';
import axiosInstance from '../../../utils/axiosInstance';

export const fetchResources = createAsyncThunk(
	'resources/fetchResources',
	async (
		{
			size,
			searchAfter,
			searchBefore,
		}: { size?: number; searchAfter?: string; searchBefore?: string },
		{ rejectWithValue }
	) => {
		try {
			const params: Record<string, any> = {};
			if (size) params.size = size;
			if (searchAfter) params.searchAfter = searchAfter;
			if (searchBefore) params.searchBefore = searchBefore;

			const response = await axiosInstance.get('/camunda/process-instance', {
				params,
			});
			return response.data.data;
		} catch (error) {
			return rejectWithValue(error);
		}
	}
);
